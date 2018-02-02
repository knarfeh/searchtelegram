import * as React from 'react';
import Helmet from 'react-helmet';
import { Link } from 'react-router';
import {
  Layout, TopBar, LayoutBody,
  LayoutResults, ActionBar, ActionBarRow, SideBar} from '../ui/layout';
import { SearchBox } from '../search';
import { SearchkitManager, SearchkitProvider } from '../../core';
// import './styles.css';

const host = "http://todo.todo"
const searchkit = new SearchkitManager(host, {useHistory: true})

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
      <SearchkitProvider searchkit={searchkit}>
        <Layout>
          <TopBar>
            <div className="my-logo">Search Telegram</div>
            <SearchBox searchOnChange={true} autofocus={true}>
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
      </SearchkitProvider>
    );
  }

}
