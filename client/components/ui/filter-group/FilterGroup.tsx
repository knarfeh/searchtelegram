import * as React from 'react';

import {
  FastClick,
  block
} from '../../../core/react';

const map = require('lodash/map')

export interface FilterGroupItemProps {
  key: string
  itemKey: string
  bemBlocks?: any
  label: string
  filter: any
  removeFilter: Function
}

export class FilterGroupItem extends React.PureComponent<FilterGroupItemProps, any> {
  constructor(props) {
    super(props)
    this.removeFilter = this.removeFilter.bind(this)
  }

  removeFilter() {
    const { removeFilter, filter } = this.props
    if (removeFilter) {
      removeFilter(filter)
    }
  }

  render() {
    const { bemBlocks, label, itemKey } = this.props
    return (
      <FastClick handler={this.removeFilter}>
        <div className={bemBlocks.item("value") } data-key={itemKey}>{label}</div>
      </FastClick>
    )
  }
}
