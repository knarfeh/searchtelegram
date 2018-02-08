import { Utils } from './../support/Utils';
import { SearchkitManager } from './../SearchkitManager';
import { ImmutableQuery } from './../query/ImmutableQuery';

const get = require('lodash/get')
const compact = require('lodash/compact')

export class Accessor {
  searchkit: SearchkitManager
  uuid: string
  results: any
  active: boolean
  translations: Object
  refCount: number

  constructor() {
    this.active = true
    this.translations = {}
    this.refCount = 0
  }

  incrementRef() {
    this.refCount++
  }

  decrementRef() {
    this.refCount--
  }

  setActive(active: boolean) {
    this.active = active
    return this
  }

  setSearchkitManager(searchkit) {
    this.searchkit = searchkit
    // this.uuid = searchkit.guid()
    console.log('WTF is guid???', this.uuid)
    this.results = this.searchkit.results
    console.log('WTF is search result', this.results)
  }

  getAggregations(defaultValue){
    const results = this.getResults()
    return results
    // return get(results, )
  }

  translate(key, interpolations?){
    let translation = (
      (this.searchkit && this.searchkit.translate(key)) ||
       this.translations[key] ||
       key)
    return Utils.translate(translation, interpolations)
  }

  getResults() {
    return this.results
  }

  setResults(results) {
    this.results = results
  }

  beforeBuildQuery() {

  }

  buildSharedQuery(query: ImmutableQuery) {
    return query
  }

  buildOwnQuery(query: ImmutableQuery) {
    return query
  }
}
