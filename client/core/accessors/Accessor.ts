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
  }

  getResults() {
    return this.results
  }

  setResults(results) {
    this.results = results
  }
}
