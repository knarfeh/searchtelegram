import { QueryAccessor } from './accessors/QueryAccessor';
import { SearchAxiosApiTransport, SearchApiTransportOptions } from './transport'
import { EventEmitter, GuidGenerator } from "./support";
import { AccessorManager } from './AccessorManager';
import { createHistoryInstance, encodeObjUrl, decodeObjString } from './history';
import { SearchRequest } from './SearchRequest';
import { ImmutableQuery } from './index';
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
  createHistory?: Function,
  getLocation?: Function,
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
  guidGenerator: GuidGenerator
  host: string
  initialLoading: boolean
  loading: boolean
  options: SearchkitOptions
  query: ImmutableQuery
  results: any
  state: any
  translateFunction: Function
  transport: SearchAxiosApiTransport

  constructor(
    host: string,
    options: SearchkitOptions = {}
  ) {
    this.options = defaults(options, {
      createHistory: createHistoryInstance,
      getLocation: ()=> typeof window !== 'undefined' && window.location,
      useHistory: true,
      searchOnload: true,
      httpHeaders: {},
    })
    console.log('constructor, options?????', this.options)
    this.host = host
    this.guidGenerator = new GuidGenerator()
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
    this.translateFunction = constant(undefined)
    this.query = new ImmutableQuery()
    this.emitter = new EventEmitter()
  }

  setupListeners() {
    this.initialLoading = !this.results
    if(this.options.useHistory) {
      console.log('usehistory??????')
      this.unlistenHistory()
      this.history = this.options.createHistory()
      this.listenToHistory()
    }
    this.runInitialSearch()
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
    this._unlistenHistory = this.history.listen((location, action)=>{
      console.log('listenTohistory')
      // action is POP when the browser modified
      if(action === "POP") {
        this._searchWhenCompleted(location)
      }
    })
  }

  _searchWhenCompleted(location){
    this.registrationCompleted.then(()=> {
      this.searchFromUrlQuery(location.search)
    }).catch((e)=> {
      console.error(e.stack)
    })
  }

  runInitialSearch() {
    if(this.options.searchOnload) {
      this.registrationCompleted.then(()=> {
        this._searchWhenCompleted(this.options.getLocation())
      })
    }
  }

  searchFromUrlQuery(query) {
    query = decodeObjString(query.replace(/^\?/, ""))
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
      console.log('performSearch, this.state', this.state)
      let url = this.options.getLocation().pathname + '?' + encodeObjUrl(this.state)
      historyMethod.call(this.history, url)
    }
  }

  getResultsAndState() {
    return {
      results: this.results,
      state: this.state
    }
  }

  _search() {
    this.state = this.accessors.getState()
    this.query = this.buildQuery()
    console.log('Done build shared query', this.query)
    const queryString = this.buildQueryString()
    this.emitter.trigger()
    this.currentSearchRequest && this.currentSearchRequest.deactivate()
    this.currentSearchRequest = new SearchRequest(
      this.transport, queryString, this
    )
    this.currentSearchRequest.run()
      .then(()=> {
        return this.getResultsAndState()
      })
  }

  translate(key) {
    return this.translateFunction(key)
  }

  buildQuery() {
    return this.accessors.buildQuery()
  }

  buildQueryString() {
    const params = this.state || this.accessors.getState()
    let keys = []
    for (let key in params) {
      if (params.hasOwnProperty(key)) {
        console.log(key, params[key]);
        if (Array.isArray(params[key])) {
          keys = keys.concat(key+'='+params[key].join(','))
        } else if(typeof params[key] === 'string') {
          keys = keys.concat(key+'='+params[key])
        }
      }
    }
    const joined = keys.join('&')
    return joined
  }

  guid(prefix){
    return this.guidGenerator.guid(prefix)
  }

  getHits() {
    return get(this.results, ["results"], [])
  }

  getHitsCount() {
    const hitsCount = get(this.results, ["total"], 0)
    console.log('hitscount', hitsCount)
    return get(this.results, ["total"], 0)
  }

  setResults(results) {
    console.log('Searchkit Manager set results', results)
    this.results = results
    this.error = null
    this.accessors.setResults(results)
    this.onResponseChange()
    // TODO
    // this.getHits()
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

