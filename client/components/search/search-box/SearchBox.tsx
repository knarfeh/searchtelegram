import * as React from 'react';
import * as PropTypes from 'prop-types';

import {
  QueryAccessor,
  SearchkitComponent,
  SearchkitComponentProps,
  ImmutableQuery
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
  throttledSearch: () => void

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
    this.throttledSearch = throttle(()=> {
      this.searchQuery(this.accessor.getQueryString())
    }, props.searchThrottleTime)
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
    // const query = this.getValue()
    // if (query.startsWith('tags:')) {
    //   console.log('WTF is getValue???', this.getValue())
    //   this.setState({input: ''})
    // }
    // console.log('value?????', this.getValue())
    // console.log('searchQuery????')
    // this.accessor.buildSharedQuery(new ImmutableQuery)
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
      return this.getAccessorValue()
    } else {
      return input
    }
  }

  getAccessorValue() {
    return (this.accessor.state.getValue() || "") + ""
  }

  onChange(e) {
    const query = e.target.value;
    console.log('onchange!!!!', query)
    if (this.props.searchOnChange) {
      this.accessor.setQueryString(query)
      this.throttledSearch()
      this.forceUpdate()
    } else {
      console.log('not search on change?????')
      this.setState( { input: query })
    }
  }

  setFocusState(focused: boolean) {
    if (!focused) {
      const { input } = this.state
      if (this.props.blurAction == 'search'
        && !isUndefined(input)
        && input != this.getAccessorValue()
      ) {
        this.searchQuery(input)
      }
      this.setState({
        focused,
        input: undefined
      })
    } else {
      this.setState({ focused })
    }
  }

  render() {
    let block = this.bemBlocks.container

    return (
      <div className={block.state({focused: this.state.focused})}>
        <form onSubmit={this.onSubmit.bind(this)}>
          <div className="sk-search-box__icon"></div>
          <input
            autoFocus={this.props.autofocus}
            className={block("text")}
            data-qa="query"
            onBlur={this.setFocusState.bind(this, false)}
            onFocus={this.setFocusState.bind(this, true)}
            onInput={this.onChange.bind(this)}
            ref="queryField"
            type="text"
            value={this.getValue()}
          />
          <input type="submit" value="search" className="sk-search-box__action" data-qa="submit"/>
        </form>
      </div>
    );
  }
}
