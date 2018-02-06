import * as React from 'react';
import { ItemComponent, CheckboxItemComponent } from './ItemComponents';
import { ListProps } from './ListProps';

const block = require('bem-cn')
const map = require('lodash/map')
const includes = require('lodash/includes')
const defaults = require('lodash/defaults')
const identity = require('lodash/identity')

export interface ItemListProps extends ListProps {
  itemComponent?: any
}

export class AbstractItemList extends React.Component<ItemListProps, {}> {
  static defaultProps: any = {
    mod: 'sk-item-list',
    showCount: true,
    itemComponent: CheckboxItemComponent,
    multiselect: true,
    selectItems: [],
  }

  isActive(option) {
    const { selectedItems, multiselect } = this.props
    if (multiselect) {
      return includes(selectedItems, option.key)
    } else {
      if (selectedItems.length == 0) return null
      return selectedItems[0] == option.key
    }
  }

  render() {
    const {
      mod, itemComponent, items, selectedItems = [], toggleItem,
      setItems, multiselect, disabled, showCount, className, docCount
    } = this.props

    const bemBlocks = {
      container: block(mod),
      option: block(`${mod}-option`)
    }

    const toggleFunc = multiselect ? toggleItem: (key => setItems([key]))

    const actions = map(items, (option) => {
      const label = option.title || option.label || option.key
      return React.createElement(itemComponent, {
        label: label,
        onClick: () => toggleFunc(option.key),
        bemBlocks: bemBlocks,
        key: option.key,
        itemKey: option.key,
        // count: 666,
        disabled: option.disabled,
        active: this.isActive(option)
      })
    })
    return (
      <div data-qa="options" className={bemBlocks.container.mix(className).state({ disabled })}>
        { actions }
      </div>
    )
  }
}

export class ItemList extends AbstractItemList {
  static defaultProps = defaults({
    itemComponent: ItemComponent
  }, AbstractItemList.defaultProps)
}

export class CheckboxItemList extends AbstractItemList {
  static defaultProps = defaults({
    itemComponent: CheckboxItemComponent
  }, AbstractItemList.defaultProps)
}

export class Toggle extends AbstractItemList {
  static defaultProps = defaults({
    itemComponent: ItemComponent,
    mod: 'sk-toggle',
    showCount: false,
  }, AbstractItemList.defaultProps)
}

export class Tabs extends AbstractItemList {
  static defaultProps = defaults({
    itemComponent: ItemComponent,
    mod: 'sk-tabs',
    showCount: false,
    multiselect: false
  }, AbstractItemList.defaultProps)
}
