import { ImmutableQuery } from './query/ImmutableQuery';
import { BaseQueryAccessor, Accessor, StatefulAccessor, noopQueryAccessor } from './accessors';
import { Utils } from './support/Utils';

const filter = require("lodash/filter")
const values = require("lodash/values")
const reduce = require("lodash/reduce")
const assign = require("lodash/assign")
const each = require("lodash/each")
const without = require("lodash/without")
const find = require("lodash/find")


type StatefulAccessors = Array<StatefulAccessor<any>>

export class AccessorManager {
  accessors: Array<Accessor>
  statefuleAccessors: {}
  queryAccessor: any

  constructor() {
    this.accessors = []
    // this.queryAccessor =
    this.statefuleAccessors = {}
  }

  getAccessors() {
    return this.accessors
  }


  getActiveAccessors() {
    return filter(this.accessors, {active: true})
  }

  getStatefuleAccessors() {
    return values(this.statefuleAccessors) as StatefulAccessors
  }

  getAccessorsByType(type) {
    return filter(this.accessors, Utils.instanceOf(type))
  }

  add(accessor) {
    console.log('accessor manager, add accessor')
    if(accessor instanceof StatefulAccessor) {
      if(accessor instanceof BaseQueryAccessor && accessor.key == 'q') {
        if(false) {

        } else {
          this.queryAccessor = accessor
        }
      }
      let existingAccessor = this.statefuleAccessors[accessor.key]
      if(existingAccessor) {
        if (existingAccessor.constructor === accessor.constructor) {
          existingAccessor.incrementRef()
          return existingAccessor
        } else {
          throw new Error(`Multiple imcompatible components with id='${accessor.key}' existing on the page`)
        }
      } else {
        this.statefuleAccessors[accessor.key] = accessor
      }
    }
    accessor.incrementRef()
    this.accessors.push(accessor)
    return accessor
  }

  remove(accessor) {
    if(!accessor) {
      return
    }
    accessor.decrementRef()
    if(accessor.refCount === 0) {
      if(accessor instanceof StatefulAccessor) {
        this.queryAccessor = noopQueryAccessor
      }
      delete this.statefuleAccessors[accessor.key]
    }
    this.accessors = without(this.accessors, accessor)
  }

  getState() {
    console.log('accessormanager get state')
    return reduce(this.getStatefuleAccessors(), (state, accessor)=> {
      return assign(state, accessor.getQueryObject())
    }, {})
  }

  setState(state) {
    console.log('set state')
    each(
      this.getStatefuleAccessors(),
      accessor=> accessor.fromQueryObject(state)
    )
  }

  notifyStateChange(oldState) {
    each(
      this.getStatefuleAccessors(),
      accessor => accessor.onStateChange(oldState)
    )
  }

  buildSharedQuery(query){
    return reduce(this.getActiveAccessors(), (query, accessor)=>{
      return accessor.buildSharedQuery(query)
    }, query)
  }

  buildOwnQuery(query){
    return reduce(this.getActiveAccessors(), (query, accessor)=>{
      return accessor.buildOwnQuery(query)
    }, query)
  }

  buildQuery(){
    each(this.getActiveAccessors(), accessor => accessor.beforeBuildQuery())
    return this.buildOwnQuery(
      this.buildSharedQuery(new ImmutableQuery())
    )
  }

  setResults(results) {
    each(this.accessors, a => a.setResults(results))
  }

  resetState() {
    each(this.getStatefuleAccessors(), a => a.resetState())
  }
}
