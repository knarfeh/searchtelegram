import * as React from 'react';
import * as PropTypes from 'prop-types';

import {
  SearchkitComponent,
  SearchkitComponentProps,
  RenderComponentType,
  renderComponent,
} from '../../../core';

const defaults = require('lodash/defaults')
const identity = require('lodash/identity')

export interface HitsStatsDisplayProps {
  bemBlocks: { container: Function }
  resultsFoundLabel: string
  // timeTaken: string | number
  hitsCount: string | number
  translate: Function
}

const HitsStatsDisplay = (props: HitsStatsDisplayProps) => {
  const { resultsFoundLabel, bemBlocks } = props
  return (
    <div className={bemBlocks.container()} data-qa="hits-stats">
      <div className={bemBlocks.container('info')} data-qa="info">
        {resultsFoundLabel}
      </div>
    </div>
  )
}

export interface HitsStatsProps extends SearchkitComponentProps {
  component?: RenderComponentType<HitsStatsDisplayProps>
  countFormatter?: (count:number)=> number | string
}

export class HitsStats extends SearchkitComponent<HitsStatsProps, any> {
  static translations: any= {
    "hitstats.results_found": "{hitCount} results found"
  }
  translations = HitsStats.translations

  static propsTypes = defaults({
    translations: SearchkitComponent.translationsPropType(
      HitsStats.translations
    ),
    countFormatter: PropTypes.func
  }, SearchkitComponent.propTypes)

  static defaultProps = {
    component: HitsStatsDisplay,
    countFormatter: identity
  }

  defineBEMBlocks() {
    return {

      container: (this.props.mod || "sk-hits-stats")
    }
  }

  render() {
    const { countFormatter } = this.props
    const hitsCount = countFormatter(this.searchkit.getHitsCount())
    const props: HitsStatsDisplayProps = {
      bemBlocks: this.bemBlocks,
      translate: this.translate,
      hitsCount: hitsCount,
      resultsFoundLabel: this.translate('hitstats.results_found', {
        hitCount: hitsCount
      })
    }
    return renderComponent(this.props.component, props)
  }
}
