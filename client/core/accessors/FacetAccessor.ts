import { FacetAccessorOptions } from './FacetAccessor';
import { FilterBasedAceessor } from './FilterBasedAccessor';
import { ArrayState } from './../state/';
// import { FieldContext } from '../query';
import { SelectedFilter } from '../query';

const assign = require("lodash/assign")
const map = require("lodash/map")
const omitBy = require("lodash/omitBy")
const isUndefined = require("lodash/isUndefined")
const keyBy = require("lodash/keyBy")
const reject = require("lodash/reject")
const each = require("lodash/each")
const identity = require("lodash/identity")

export interface FacetAccessorOptions {
  operator?:string
  title?:string
  id?:string
  size:number
  facetsPerPage?:number
  translations?:Object
  include?:Array<string> | string
  exclude?:Array<string> | string
  orderKey?:string
  orderDirection?:string
  min_doc_count?:number
  loadAggregations?: boolean
  // fieldOptions?:FieldOptions
}

export interface ISizeOption {
  label: string
  size: number
}

export class FacetAccessor extends FilterBasedAceessor<ArrayState> {
  state = new ArrayState()
  options: any
  defaultSize: number
  size: number
  uuid: string
  // fieldContext: FieldContext

  constructor(key, options: FacetAccessorOptions) {
    super(key, options.id)
    this.options = options
  }

  getRawBuckets() {
    return this.getAggregations('')
  }

  getBuckets() {
    let rawBuckets: Array<any> = this.getRawBuckets()
    let keyIndex = {}
    each(rawBuckets, (item) => {
      item.key = String(item.key)
      // keyIndex[item.key]
      item['title'] = item.key
    })
    return rawBuckets
  }

  buildSharedQuery(query){
    var filters = this.state.getValue()
    console.log('WTF is buildSharedQuery??? filter', filters)
    console.log('And wtf is query', query)
    var selectedFilters:Array<SelectedFilter> = map(filters, (filter)=> {
      console.log('WTF is option???', this.options)
      return {
        name:this.options.title || this.translate(this.options.field),
        value:this.translate(filter),
        id:this.options.id,
        remove:()=> this.state = this.state.remove(filter)
      }
    })
    query = query.addSelectedFilters(selectedFilters)
    console.log('buildShareQuery, wtf is query', query)
    return query
  }
}
