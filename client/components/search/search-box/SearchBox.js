import * as React from 'react';
import * as PropTypes from 'prop-types';

import { SearchkitComponentProps } from 'searchkit';

const defaults = require('lodash/defaults');
const throttle = require('lodash/throttle');
const assign = require('lodash/assign');
const isUndefined = require('lodash/isUndefined');

export class SearchBox extends Component {
  render() {
    return (
      <div className="searchboxclassname">
      <form>
        <div className="TODO"></div>
        <input type="text"
        data-qa="query"
        className="TODO"
        placeholder="search"
        ref="queryField"/>
        </form>
      </div>
    );
  }
}
