import { BaseQueryAccessor } from './BaseQueryAccessor';
import { SelectedFilter } from '../query';

const map = require('lodash/map')

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
    super(key)
    this.options = options
    // this.options.queryFields = this.options.queryFields || ["_all"]
  }

  fromQueryObject(ob) {
    super.fromQueryObject(ob)
  }



  buildSharedQuery(query) {
    // let filters = this.state.getValue()
    // console.log('!!!query accesor build shared query, filters???', filters)
    // var tagsString
    // if (filters !== null) {
    //   tagsString = filters.toString().split(':')
    //   console.log('tagsString???', tagsString)
    // } else {
    //   tagsString = []
    // }

    // console.log('WTF is tagsString length???', tagsString.length)
    // if (tagsString.length === 2) {
    //   console.log('Got in')
    //   var selectedFilters:Array<SelectedFilter> = map(filters, (filter)=> {
    //     return {
    //       name: "Tags",
    //       value: tagsString[1],
    //       id: "id",
    //       remove:()=> this.state = this.state.clear()
    //     }
    //   })
    //   query = query.addSelectedFilters(selectedFilters)
    // }

    return query
  }

}
