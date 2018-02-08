const update = require('immutability-helper')
import { Utils } from '../support/Utils';
import { SelectedFilter } from './SelectedFilter';

const omitBy = require('lodash/omitBy')
const omit = require('lodash/omit')
const values = require('lodash/values')
const pick = require('lodash/pick')
const merge = require('lodash/merge')
const isUndefined = require('lodash/isUndefined')

export type SourceFilterType = string | Array<string> | boolean

export class ImmutableQuery {
  index: any
  query: any
  static defautlIndex: any = {
    queryString: "",
    filtersMap: {},
    selectedFilters: [],
    queries: [],
    filters: [],
    size: 0
  }

  constructor(index = ImmutableQuery.defautlIndex) {
    this.index = index
  }

  hasFilters() {
    return this.index.filters.length > 0
  }

  // hasFiltersOrQuery () {
    // return (this.index.q)
  // }
  addSelectedFilter(selectedFilter: SelectedFilter) {
    return this.addSelectedFilters([selectedFilter])
  }

  addSelectedFilters(selectedFilters: Array<SelectedFilter>) {
    return this.update({
      selectedFilters: {$push: selectedFilters}
    })
  }

  getSelectedFilters(): Array<SelectedFilter> {
    return this.index.selectedFilters
  }

  update(updateDef) {
    return new ImmutableQuery(update(this.index, updateDef))
  }

  addFilter(key, filter) {
    return this.update({
      filters: { $push: [filter]},
      filtersMap: { $merge: { [key]: filter }}
    })
  }

  getFilters(keys=[]) {
    // return this.getFiltersWithoutKeys(keys)
  }
}
