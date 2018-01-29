import { SearchAxiosApiTransport, SearchApiTransportOptions } from './transport'
import { EventEmitter } from "./support";
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
  host: string
  state: any
  currentSearchRequest: Function
  history
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
    this.emitter = new EventEmitter()
  }

  setupListerners() {
    this.initialLoading = true
    if(this.options.useHistory) {

    } else {
      // this.runInitialSearch()
    }
  }
}

