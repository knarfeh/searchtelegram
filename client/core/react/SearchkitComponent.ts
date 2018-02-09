import * as React from "react"
import * as PropTypes from 'prop-types';
import { SearchkitManager } from "../SearchkitManager";
import { Accessor } from "../accessors/Accessor"
import { Utils } from "../support"
import { ImmutableQuery } from './../query';
var block = require('bem-cn');
const keys = require("lodash/keys")
const without = require("lodash/without")
const transform = require("lodash/transform")
const mapValues = require("lodash/mapValues")

export interface SearchkitComponentProps {
  mod?: string
  className?: string
  translations?: Object
  searchkit?: SearchkitManager
  key?: string
}

export class SearchkitComponent<P extends SearchkitComponentProps,S> extends React.Component<P,S> {
  searchkit: SearchkitManager
  accessor: Accessor
  stateListenerUnsubscribe: Function
  translations: Object = {}
  unmounted = false

  static contextTypes: React.ValidationMap<any> = {
    searchkit: React.PropTypes.instanceOf(SearchkitManager)
  }

  static translationsPropType = (translations)=> {
    return PropTypes.shape(mapValues(translations, ()=> PropTypes.string))
  }

  static propTypes: any = {
    mod: React.PropTypes.string,
    className: React.PropTypes.string,
    searchkit: React.PropTypes.instanceOf(SearchkitManager)
  }

  constructor(props?) {
    super(props)
  }

  defineBEMBlocks() {
    return null;
  }

  defineAccessor(): Accessor{
    return null;
  }

  translate(key, interpolations?) {
    let translation = (
      (this.searchkit.translate(key)) ||
      (this.props.translations && this.props.translations[key]) ||
      this.translations[key] || key)
    return Utils.translate(translation, interpolations)
  }

  get bemBlocks() {
    return transform(this.defineBEMBlocks(), (result: any, cssClass, name)=> {
      result[name] = block(cssClass);
    })
  }

  _getSearchkit() {
    return this.props.searchkit || this.context['searchkit']
  }

  componentWillMount() {
    this.searchkit = this._getSearchkit()
    if(this.searchkit) {
      this.accessor = this.defineAccessor()
      if(this.accessor) {
        this.accessor = this.searchkit.addAccessor(this.accessor)
      }
      this.stateListenerUnsubscribe = this.searchkit.emitter.addListener(()=> {
        if(!this.unmounted) {
          this.forceUpdate();
        }
      })
    } else {
      console.warn("No searchkit found in props or context for " + this.constructor["name"])
    }
  }

  componentWillUnMount() {
    if(this.stateListenerUnsubscribe) {
      this.stateListenerUnsubscribe()
    }
    if(this.searchkit && this.accessor) {
      this.searchkit.removeAccessor(this.accessor)
    }
    this.unmounted = true
  }

  getQuery(): ImmutableQuery {
    return this.searchkit.query
  }

  getHits() {
    return this.searchkit.getHits()
  }

  isInitialLoading() {
    return this.searchkit.initialLoading
  }

}
