import * as React from "react";
import * as PropTypes from "prop-types";

import {
	SearchkitComponent,
	// PageSizeAccessor,
	// HighlightAccessor,
	// CustomHighlightAccessor,
	SearchkitComponentProps,
	SourceFilterType,
	// SourceFilterAccessor,
	HitsAccessor,
	RenderComponentType,
	RenderComponentPropType,
	renderComponent,
	block
} from "../../../core";

const map = require("lodash/map")
const defaults = require("lodash/defaults")

export interface HitItemProps {
  key: string,
  bemBlocks?: any,
  results: any
}

export class HitItem extends React.PureComponent<HitItemProps, any> {
  render() {
    return (
      <div data-qa="hit"
        className={this.props.bemBlocks.item().mix(this.props.bemBlocks.container("item"))}>
        {this.props.results._id}
      </div>
    )
  }
}

export interface HitsListProps {
  mod?: string,
  className?: string,
  itemComponent?: RenderComponentType<HitItemProps>,
  hits: Array<Object>
}

export class HitsList extends React.PureComponent<HitsListProps, any> {
  static defaultProps = {
    mod: "sk-hits",
    itemComponent: HitItem
  }

  static propTypes = {
    mod: PropTypes.string,
    className: PropTypes.string,
    itemComponent: RenderComponentPropType,
    hits: PropTypes.any
  }

  render() {
    const { hits, mod, className, itemComponent } = this.props
    const bemBlocks = {
      container: block(mod).el,
      item: block(`${mod}-hit`).el
    }
    return (
      <div data-qa="hits" className={bemBlocks.container().mix(className)}>
        {map(hits, (result: any, index) => {
          return renderComponent(itemComponent, {
            key: result._id, result, bemBlocks, index
          })
        })}
      </div>
    )
  }
}

export interface HitsProps extends SearchkitComponentProps {
  hitsPerPage?: number
  itemComponent?: RenderComponentType<HitItemProps>
  listComponent?: RenderComponentType<HitsListProps>
  scrollTo?: boolean | string
}

export class Hits extends SearchkitComponent<HitsProps, any> {
  hitsAccessor: HitsAccessor
  static propTypes = defaults({
    hitsPerpage: PropTypes.number,
    itemComponent: RenderComponentPropType,
    listComponent: RenderComponentPropType
  }, SearchkitComponent.propTypes)

  public static defaultsProps = {
    listComponent: HitsList,
    scrollTo: "body"
  }

  componentWillMount() {
    super.componentWillMount()
    this.hitsAccessor = new HitsAccessor({ scrollTo: this.props.scrollTo })
    this.searchkit.addAccessor(this.hitsAccessor)
  }

  render() {
    let hits: Array<Object> = this.getHits()
    let hasHits = hits.length > 0

    if (!this.isInitialLoading() && hasHits) {
      const { listComponent, mod, className, itemComponent } = this.props
      return renderComponent(listComponent, {
        hits, mod, className, itemComponent
      })
    }
    return null
  }
}
