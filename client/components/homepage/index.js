import React, { Component } from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import {
  Layout, TopBar, LayoutBody,
  LayoutResults, ActionBar, ActionBarRow, SideBar} from '../ui/layout';

export default class Homepage extends Component {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */

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
            <div className="st-logo">Telegram Search</div>
            {/* <SearchBox autofocus={true} searchOnChange={true} prefixQueryFields={["actors^1","type^2","languages","title^10"]}/> */}
          </TopBar>
          <LayoutBody>
            <SideBar>
              <p>side bar</p>
            </SideBar>
          </LayoutBody>
        </Layout>
      </div>
    );
  }

}
