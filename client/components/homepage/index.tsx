import * as React from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import {
  Layout, TopBar, LayoutBody,
  LayoutResults, ActionBar, ActionBarRow, SideBar} from '../ui/layout';
import { SearchBox } from '../search';
// import './styles.css';

export default class Homepage extends React.Component {
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
            <div className="my-logo">Search Telegram</div>
            <SearchBox>
            </SearchBox>
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
