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

  setFilter(keys) {
    this.accessor.state = this.accessor.state.setValue(keys)
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
    console.log('TODO, return items')
    const items = [
      {
        key: 'key1',
        title: 'title1'
      },
      {
        key: 'key2',
        title: 'title2'
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
