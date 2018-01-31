import { State } from '../state';
import { Accessor } from './Accessor';

export class StatefulAccessor<T extends State<any>> extends Accessor {
  key: string
  urlKey: string
  state: T
  resultsState: T

  constructor(key, urlString?) {
    super()
    this.key = key
    this.uuid = this.key + this.uuid
    this.urlKey = urlString || key && key.replace(/\./g, "_")
    this.urlWithState = this.urlWithState.bind(this)
  }

  urlWithState(state: T) {
    return this.searchkit.buildSearchUrl({ [this.urlKey]: state})
  }
}
