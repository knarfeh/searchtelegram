import { Accessor } from './Accessor'

export interface HitsOptions {
  scrollTo: string | boolean
}

export class HitsAccessor extends Accessor {
  constructor(public options: HitsOptions) {
    super()
  }

  setResults(results) {
    super.setResults(results)
    // this.scro
  }

  scrollIfNeeded() {
    // if(this.searchkit.)
  }

  getScrollSelector() {
    return (this.options.scrollTo == true) ? "body": this.options.scrollTo.toString();
  }
}
