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
    // console.log('StatefulAccessor?????', this.urlWithState)
  }

  urlWithState(state: T) {
    console.log('url with state!!!!!, this.urlKey', this.urlKey)
    return this.searchkit.buildSearchUrl({ [this.urlKey]: state})
  }

  fromQueryObject(ob) {
    console.log('Statefulaccessor, fromQueryObject????, this.urlKey???', this.urlKey)
    let value = ob[this.urlKey]
    console.log('WTF is value: ', value)
    this.state = this.state.setValue(value)
  }

  getQueryObject() {
    let val = this.state.getValue()
    return (val) ? {
      [this.urlKey]: val
    } : {}
  }

  setSearchkitManager(searchkit){
    super.setSearchkitManager(searchkit)
    // this.setResultsState()
  }

  onStateChange(oldState) {
    console.log('stateful accessor, onStateChange???')
  }
}
