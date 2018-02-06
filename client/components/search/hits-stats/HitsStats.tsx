import * as React from 'react';
import * as PropTypes from 'prop-types';

import {
  SearchkitComponent,
  SearchkitComponentProps,
  renderComponent,
} from '../../../core';

const defaults = require('lodash/defaults')
const identify = require('lodash/identify')

export interface HitsStatsDisplayProps {
  bemBlocks: { container: Function }
  resultsFoundLabel: string
  timeTaken: string | number
  hitsCount: string | number
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

// export interface
