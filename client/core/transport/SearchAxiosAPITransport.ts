import axios from "axios";
import { AxiosInstance } from 'axios';
import { SearchAPITransport } from "./SearchAPITransport";

const defaults = require("lodash/defaults")

export interface SearchApiTransportOptions {
  headers?: object,
  searchUrlPath?: string
};

export class SearchAxiosApiTransport extends SearchAPITransport {
  static timeout = 5000;
  axios: AxiosInstance
  options: SearchApiTransportOptions;

  constructor(public host: string, options: SearchApiTransportOptions={}) {
    super();
    this.options = defaults(options, {
      headers: {},
      searchUrlPath: "/search"
    });

    this.axios = axios.create({
      baseURL: this.host,
      timeout: SearchAxiosApiTransport.timeout,
      headers: this.options.headers
    })
  }


  search(queryString: String) {
    return this.axios.get(this.options.searchUrlPath+'?'+queryString+'&page=1&page_size=100')
      .then(this.getData)
  }

  getData(response) {
    return response.data
  }
};
