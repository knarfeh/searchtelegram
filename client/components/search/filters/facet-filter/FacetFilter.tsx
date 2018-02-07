import * as React from 'react';

import { FacetFilterProps, FacetFilterPropTypes } from './FacetFilterProps';
import {
  FacetAccessor, SearchkitComponent, ISizeOption,
  renderComponent,
 } from '../../../../core';

import { CheckboxItemList, Panel } from '../../../ui';

const defaults = require('lodash/defaults')
const identity = require('lodash/identity')

export class FacetFilter<T extends FacetFilterProps> extends SearchkitComponent<T, any> {
  accessor: FacetAccessor

  static propTypes = FacetFilterPropTypes

  static defaultProps = {
    listComponent: CheckboxItemList,
    containerComponent: Panel,
    size: 50,
    collapsable: false,
    showCount: true,
    showMore: true,
  }

  constructor(props) {
    super(props)
    // this.toggleView
  }

  componentDidMount() {
    var self = this;
}

  getAccessorOptions() {
    const {
      field, id, title, size, orderKey, orderDirection,
    } = this.props
    return {
      id, title, size, field, orderKey, orderDirection,
    }
  }

  defineAccessor() {
    return new FacetAccessor(
      this.props.field, this.getAccessorOptions()
    )
  }

  defineBEMBlocks() {
    var blockName = this.props.mod || 'sk-refinement-list'
    return {
      container: blockName,
      option: `${blockName}-option`
    }
  }

  toggleFilter(key) {
    console.log('toggle filter????', key)
    this.accessor.state = this.accessor.state.toggle(key)
    this.searchkit.performSearch()
  }

  setFilter(keys) {
    console.log('Facet set filter, keys???', keys)
    this.accessor.state = this.accessor.state.setValue(keys)
    this.searchkit.performSearch()
  }

  toggleViewMoreOption(option: ISizeOption) {
    console.log('TODO: toggle view more option')
    this.searchkit.performSearch()
    // this.accessor.setV
  }

  hasOptions(): boolean {
    console.log('TODO, hasOptions')
    return true
  }

  getSelectedItems() {
    return this.accessor.state.getValue()
  }

  getItems() {
    // this.accessor.getAggregations()
    // const result = this.accessor.getResults()
    // console.log('result from accessor???', result)
    // Object.keys(result).forEach(function(key) {
    //   console.log(key, result[key]);
    //   result['title'] = result[key]
    // });
    // const result = this.accessor.getBuckets()
    // console.log('now, result: ', result)
    // return result
    // console.log('get items, test????', test)
    // console.log('TODO, return items')
    const items = [
      {
        key: 'people',
        title: 'people'
      },
      {
        key: 'channel',
        title: 'channel'
      },
      {
        key: 'group',
        title: 'group'
      }
    ]
    return items
  }

  render() {
    const { listComponent, containerComponent, showCount, title, id } = this.props
    return renderComponent(containerComponent, {
      title,
      className: id ? `filter--${id}` : undefined,
      disabled: !this.hasOptions()
    }, [
      renderComponent(listComponent, {
        toggleItem: this.toggleFilter.bind(this),
        key: 'listComponent',
        items: this.getItems(),
        itemComponent: this.props.itemComponent,
        selectedItems: this.getSelectedItems(),
        setItems: this.setFilter.bind(this),
        docCount: 'TODO: getDocCount',
        showCount,
      }),
      // this.renderShowMore()
    ]);
  }
}
