import { SearchAxiosApiTransport } from "./transport";
import { SearchkitManager } from "./SearchkitManager";

export class SearchRequest {
  active: boolean

  constructor(
    public transport: SearchAxiosApiTransport,
    // public query: Object,
    public queryString: String,
    public searchkit: SearchkitManager
  ) {
    this.active = true;
  }

  run() {
    console.log('todo: bind result')
    return this.transport.search(this.queryString).then(
      this.setResults.bind(this)
    ).catch (
      this.setError.bind(this)
    )
  }

  deactivate() {
    this.active = false
  }

  setResults(results) {
    console.log('result?????', results)
    if(this.active) {
      this.searchkit.setResults(results)
    }
  }

  setError(error) {
    if(this.active) {
      // this.searchkit.setError(error)
    }
  }
}
