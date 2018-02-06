import { Utils } from './../support/Utils';
import { SearchkitManager } from './../SearchkitManager';

const get = require('lodash/get')
const compact = require('lodash/compact')

export class Accessor {
  searchkit: SearchkitManager
  uuid: string
  results: any
  active: boolean
  refCount: number

  constructor() {
    this.uuid = Utils.guid()
    this.active = true
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

  getResults() {
    return this.results
  }

  setResults(results) {
    this.results = results
  }
}
