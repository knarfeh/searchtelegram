import { BaseQueryAccessor } from './BaseQueryAccessor';

const assign = require('lodash/assign');

export interface SearchOptions {
  queryFields?:Array<string>
  queryOptions?:any
  prefixQueryFields?:Array<string>
  prefixQueryOptions?:Object
  title?: string
  addToFilters?:boolean
  queryBuilder?:Function
  onQueryStateChange?:Function
}

export class QueryAccessor extends BaseQueryAccessor {
  options: SearchOptions

  constructor(key, options={}) {
    console.log('queryaccessor construct')
    super(key)
    this.options = options
    // this.options.queryFields = this.options.queryFields || ["_all"]
  }


}
