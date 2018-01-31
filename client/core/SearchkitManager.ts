import { QueryAccessor } from './accessors/QueryAccessor';
import { SearchAxiosApiTransport, SearchApiTransportOptions } from './transport'
import { EventEmitter } from "./support";
import { AccessorManager } from './AccessorManager';
const defaults = require("lodash/defaults")
const constant = require("lodash/constant")
const identity = require("lodash/identity")
const map = require("lodash/map")
const isEqual = require("lodash/isEqual")
const get = require("lodash/get")
const after = require("lodash/after")
require('es6-promise').polyfill()

export interface SearchkitOptions {
  useHistory?: boolean,
  searchOnLoad?: boolean,
  httpHeaders?: Object,
  basicAuth?: string,
  transport?: SearchAxiosApiTransport,
  searchUrlPath?: string
}

export class SearchkitManager {
  private registrationCompleted: Promise<any>
  accessors: AccessorManager
  host: string
  state: any
  currentSearchRequest: Function
  history
  query
  options: SearchkitOptions
  transport: SearchAxiosApiTransport
  emitter: EventEmitter
  initialLoading: boolean

  constructor(
    host: string,
    options: SearchkitOptions = {}
  ) {
    this.options = defaults(options, {
      useHistory: true,
      httpHeaders: {},
      searchOnload: true
    })
    this.host = host
    this.transport = this.options.transport || new SearchAxiosApiTransport(host, {
      headers: this.options.httpHeaders,
      searchUrlPath: this.options.searchUrlPath
    })
    this.accessors = new AccessorManager()
    this.emitter = new EventEmitter()
  }

  setupListeners() {
    this.initialLoading = true
    if(this.options.useHistory) {
      // this.un
    } else {
      this.runInitialSearch()
    }
  }

  addAccessor(accessor) {
    accessor.setSearchkitManager(this)
    return this.accessors.add(accessor)
  }

  removeAccessor(accessor) {
    this.accessors.remove(accessor)
  }

  resetState() {
    this.accessors.resetState()
  }

  runInitialSearch() {
    if(this.options.searchOnLoad) {
      this.registrationCompleted.then(()=> {
        // this._search()
      })
    }
  }



  buildSearchUrl(extraParams = {}) {
    const params = defaults(extraParams, this.state || this.accessors.getState())
  }

  reloadSearch() {
    // delete this.query
    this.performSearch()
  }

  performSearch(replaceState=false, notifyState=true) {
    if(notifyState && !isEqual(this.accessors.getState(), this.state)) {
      this.accessors.notifyStateChange(this.state)
    }
    // this._search()
  }

  _search() {
    this.state = this.accessors.getState()
    console.log('this query???', this.query.getJSON())
  }

}

