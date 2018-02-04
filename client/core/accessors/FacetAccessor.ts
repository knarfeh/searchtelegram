import { FacetAccessorOptions } from './FacetAccessor';
import { FilterBasedAceessor } from './FilterBasedAccessor';
import { ArrayState } from './../state/';

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

  constructor(key, options: FacetAccessorOptions) {
    super(key, options.id)
  }
}
