import * as React from 'react';
import * as PropTypes from 'prop-types';

import {
  QueryAccessor,
  SearchkitComponent,
  SearchkitComponentProps
} from "../../../core"

const defaults = require('lodash/defaults');
const throttle = require('lodash/throttle');
const assign = require('lodash/assign');
const isUndefined = require('lodash/isUndefined');

export interface SearchBoxProps extends SearchkitComponentProps {
  searchOnChange?: boolean
  searchThrottleTime?: number
  queryFields?: Array<string>
  queryOptions?: any
  autofocus?: boolean
  id?: string
  mod?: string
  placeholder?: string
  blurAction?: "search" | "restore"
}

export class SearchBox extends SearchkitComponent<SearchBoxProps, any> {
  accessor: QueryAccessor
  lastSearchMs: number

  static defaultProps = {
    id: 'q',
    mod: 'sk-search-box',
    searchThrottleTime: 200,
    blurAction: 'search'
  }

  static propTypes = defaults({
    id: React.PropTypes.string,
    searchOnChange: React.PropTypes.bool,
    searchThrottleTime: React.PropTypes.number,
    mod: React.PropTypes.string,
    placeholder: React.PropTypes.string,
  }, SearchkitComponent.propTypes)

  constructor (props: SearchBoxProps) {
    super(props);
    this.state = {
      focused: false,
      input: undefined
    }
    this.lastSearchMs = 0

  }

  defineBEMBlocks() {
    return { container: this.props.mod };
  }

  defineAccessor() {
    const {
      id, searchOnChange, queryOptions,
    } = this.props
    return new QueryAccessor(id, {
      queryOptions: assign({}, queryOptions),
      onQueryStateChange: () => {
        if(!this.unmounted && this.state.input) {
          this.setState({input: undefined})
        }
      }
    })
  }

  onSubmit(event) {
    event.preventDefault()
    console.log('value?????', this.getValue())
    console.log('searchQuery????')
    this.searchQuery(this.getValue())
  }

  searchQuery(query) {
    let shouldResetOtherState = false
    this.accessor.setQueryString(query, shouldResetOtherState)
    let now = +new Date
    let newSearch = now - this.lastSearchMs <= 2000
    this.lastSearchMs = now
    this.searchkit.performSearch(newSearch)
  }

  getValue() {
    const { input } = this.state
    if (isUndefined(input)) {
      // return this.getAc
    } else {
      return input
    }
  }

  getAccessorValue() {
    return (this.accessor.state.getValue() || "") + ""
  }

  onChange(e) {
    const query = e.target.value;
    if (this.props.searchOnChange) {

    } else {
      this.setState( { input: query })
    }
  }

  render() {
    return (
      <div className="sk-search-box">
        <form onSubmit={this.onSubmit.bind(this)}>
          <div className="sk-search-box__icon"></div>
          <input
            type="text"
            value={this.getValue()}
            onInput={this.onChange.bind(this)}
            className="sk-search-box__text"
          />
          <input type="submit" value="search" className="sk-search-box__action" data-qa="submit"/>
        </form>
      </div>
    );
  }
}
