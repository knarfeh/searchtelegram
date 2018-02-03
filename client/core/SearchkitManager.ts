import { QueryAccessor } from './accessors/QueryAccessor';
import { SearchAxiosApiTransport, SearchApiTransportOptions } from './transport'
import { EventEmitter } from "./support";
import { AccessorManager } from './AccessorManager';
import { createHistory } from './history';
import { SearchRequest } from './SearchRequest';
// import * as Promise from 'bluebird';

const defaults = require("lodash/defaults")
const constant = require("lodash/constant")
const identity = require("lodash/identity")
const map = require("lodash/map")
const isEqual = require("lodash/isEqual")
const get = require("lodash/get")
const after = require("lodash/after")
const qs = require('qs')

require('es6-promise').polyfill()

export interface SearchkitOptions {
  useHistory?: boolean,
  searchOnload?: boolean,
  httpHeaders?: Object,
  basicAuth?: string,
  transport?: SearchAxiosApiTransport,
  searchUrlPath?: string
}

export class SearchkitManager {
  private registrationCompleted: any
  _unlistenHistory: Function
  accessors: AccessorManager
  completeRegistration: Function
  currentSearchRequest: SearchRequest
  error: any
  emitter: EventEmitter
  history
  host: string
  initialLoading: boolean
  loading: boolean
  options: SearchkitOptions
  query
  results: any
  state: any
  transport: SearchAxiosApiTransport

  constructor(
    host: string,
    options: SearchkitOptions = {}
  ) {
    this.options = defaults(options, {
      useHistory: true,
      searchOnload: true,
      httpHeaders: {},
    })
    console.log('constructor, options?????', this.options)
    this.host = host
    this.transport = this.options.transport || new SearchAxiosApiTransport(host, {
      headers: this.options.httpHeaders,
      searchUrlPath: this.options.searchUrlPath
    })
    this.accessors = new AccessorManager()
    // 'Promise' only refers to a type, but is being used as a value here.
    // https://github.com/ReactiveX/rxjs/issues/2422
    this.registrationCompleted = new Promise((resolve)=> {
      this.completeRegistration = resolve
    })
    this.emitter = new EventEmitter()
  }

  setupListeners() {
    this.initialLoading = true
    if(this.options.useHistory) {
      console.log('usehistory??????')
      this.unlistenHistory()
      this.history = createHistory()
      this.listenToHistory()
      // this.un
    } else {
      console.log('setuplisteners, not use history')
      this.runInitialSearch()
    }
  }

  addAccessor(accessor) {
    console.log('searchkitmanager add accessor')
    accessor.setSearchkitManager(this)
    return this.accessors.add(accessor)
  }

  removeAccessor(accessor) {
    this.accessors.remove(accessor)
  }

  resetState() {
    this.accessors.resetState()
  }

  listenToHistory() {
    console.log('this.options?????', this.options.searchOnload)
    let callsBeforeListen = (this.options.searchOnload) ? 1: 2
    console.log('callsbeforelisten????', callsBeforeListen)
    this._unlistenHistory = this.history.listen(after(callsBeforeListen, (location)=>{
      console.log('listenTohistory')
      // action is POP when the browser modified
      if(location.action === "POP") {
        this.registrationCompleted.then(()=> {
          this.searchFromUrlQuery(location.query)
        }).catch((e)=> {
          console.error(e.stack)
        })
      }
    }))
  }

  runInitialSearch() {
    if(this.options.searchOnload) {
      this.registrationCompleted.then(()=> {
        this._search()
      })
    }
  }

  searchFromUrlQuery(query) {
    console.log('search from url query')
    this.accessors.setState(query)
    this._search()
  }

  unlistenHistory() {
    console.log('unlistenHistory, todo')
  }


  buildSearchUrl(extraParams = {}) {
    const params = defaults(extraParams, this.state || this.accessors.getState())
    const queryString = qs.stringify(params, { encode: true})
    return window.location.pathname + '?' + queryString
  }

  reloadSearch() {
    // delete this.query
    this.performSearch()
  }

  performSearch(replaceState=false, notifyState=true) {
    if(notifyState && !isEqual(this.accessors.getState(), this.state)) {
      this.accessors.notifyStateChange(this.state)
    }
    this._search()
    if(this.options.useHistory) {
      const historyMethod = (replaceState) ?
      this.history.replace : this.history.push
      historyMethod({pathname: window.location.pathname, query: this.state})
    }
  }

  _search() {
    this.state = this.accessors.getState()
    const params = this.state || this.accessors.getState()
    const queryString = qs.stringify(params, { encode: true})
    // this.query = this.buildQuery()
    this.emitter.trigger()
    this.currentSearchRequest && this.currentSearchRequest.deactivate()
    this.currentSearchRequest = new SearchRequest(
      this.transport, queryString, this
    )
    this.currentSearchRequest.run()
  }

  setResults(results) {
    console.log('got results', results)
    this.results = results
    this.error = null
    this.accessors.setResults(results)
    // this.onResponseChange()
    // TODO
  }

  setError(error) {
    this.error = error
    console.error(this.error)
    this.results = null
    this.accessors.setResults(null)
    this.onResponseChange()
  }

  onResponseChange() {
    this.loading = false
    this.initialLoading = false
    this.emitter.trigger()
  }

  compareResults(priviousResults, results) {
    console.log('TODO: compare results')
  }
}

