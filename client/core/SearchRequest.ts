import { SearchAxiosApiTransport } from "./transport";
import { SearchkitManager } from "./SearchkitManager";

export class SearchRequest {
  active: boolean

  constructor(
    public transport: SearchAxiosApiTransport,
    public query: Object,
    public searchkit: SearchkitManager
  ) {
    this.active = true;
  }

  run() {
    return this.transport.search(this.query).then(
    //   this.set
    )
  }
}