import * as React from 'react';
import {
  FastClick
} from '../../../core';

export interface NoHitsDisplayProps {
  noResultsLabel: string
  resetFiltersFn: Function
  translate: Function
  bemBlocks: { container: Function }
  query: string
  filtersCount: number
}

export class NoHitsDisplay extends React.Component<NoHitsDisplayProps, any> {

  getSuggestionAction() {
    // TODO
    return null
  }

  getResetFilterAction() {
    const { filtersCount, query, resetFiltersFn, bemBlocks, translate } = this.props

    if (filtersCount > 0) {
      return (
        <FastClick handler={resetFiltersFn}>
          <div className={bemBlocks.container('step-action')}>
          </div>
        </FastClick>
      )
    }
  }

  render() {
    const { bemBlocks, noResultsLabel } = this.props
    return (
      <div data-qa="no-hits" className={bemBlocks.container()}>
        <div className={bemBlocks.container("info")}>
          {noResultsLabel}
        </div>
        <div className={bemBlocks.container("steps")}>
          {/* {this.get} */}
        </div>
      </div>
    )
  }
}
