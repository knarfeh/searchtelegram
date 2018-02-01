import { StatefulAccessor } from './StatefulAccessor';
import { ValueState } from './../state/ValueState';

export class BaseQueryAccessor extends StatefulAccessor<ValueState> {

  constructor(key) {
    super(key)
    this.state = new ValueState()
  }

  keepOnlyQueryState() {
    this.setQueryString(this.getQueryString(), true)
  }

  setQueryString(queryString, withReset=false) {
    console.log('setQueryString!!!!, queryString???', queryString)
    if (withReset) {
      this.searchkit.resetState()
    }
    this.state = this.state.setValue(queryString)
    // this.urlWithState(this.state) // TODELETE
  }

  getQueryString() {
    return this.state.getValue()
  }
}

export class NoopQueryAccessor extends BaseQueryAccessor {
  keepOnlyQueryState(){
    console.warn("keepOnlyQueryState called, No Query Accessor exists")
  }

  setQueryString(queryString, withReset=false){
    console.warn("setQueryString called, No Query Accessor exists")
  }

  getQueryString(){
    console.warn("getQueryString called, No Query Accessor exists")
    return ""
  }
}

export const noopQueryAccessor = new NoopQueryAccessor(null)
