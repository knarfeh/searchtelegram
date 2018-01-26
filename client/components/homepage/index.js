import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import {
  Layout, TopBar, LayoutBody,
  LayoutResults, ActionBar, ActionBarRow, SideBar} from '../ui/layout';
import './styles.css';

export default class Homepage extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */
  constructor(props) {
    super(props);
    console.log(window.location);
    console.log('constructor');
  }

  render() {
    return (
      <div>
        <Helmet
          title='Search Telegram'
          meta={[
            {
              property: 'og:title',
              content: 'Search Engine for Telegram'
            }
          ]} />
        <Layout>
          <TopBar>
            <div className="my-logo">Telegram Search</div>
            <div className="sk-search-box">
              <form>
                <div className="sk-search-box__icon"></div>
                <input type="text"
                data-qa="query"
                className="TODO"
                placeholder="search"
                ref="queryField"/>
              </form>
            </div>
          </TopBar>
          <LayoutBody>
            <SideBar>
              <p>side bar</p>
            </SideBar>
            <LayoutResults>
              <p>Result</p>
            </LayoutResults>
          </LayoutBody>
        </Layout>
      </div>
    );
  }

}
