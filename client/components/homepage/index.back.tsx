// import * as React from 'react';
// import Helmet from 'react-helmet';
// import { Link } from 'react-router';
// import {
//   Layout, TopBar, LayoutBody,
//   LayoutResults, ActionBar, ActionBarRow, SideBar} from '../ui/layout';
// import {
//   SearchBox,
//   RefinementListFilter
// } from '../search';
// import { SearchkitManager, SearchkitProvider} from '../../core';
// import { createHistoryInstance } from '../../core/history';
// import { HitsStats } from '../search/hits-stats/HitsStats';
// import { GroupedSelectedFilters } from '../search';
// // import './styles.css';

// const host = "/api/v1"
// const searchkit = new SearchkitManager(host, {useHistory: true, createHistory: createHistoryInstance})

// export default class Homepage extends React.Component {
//   /*eslint-disable */
//   static onEnter({store, nextState, replaceState, callback}) {
//     // Load here any data.
//     callback(); // this call is important, don't forget it
//   }
//   /*eslint-enable */
//   constructor(props) {
//     super(props);
//     console.log(window.location);
//     console.log('constructor');
//   }

//   render() {
//     return (
//       <SearchkitProvider searchkit={searchkit}>
//         <Layout>
//           <TopBar>
//             <div className="my-logo">Search Telegram</div>
//             <SearchBox searchOnChange={true} autofocus={true}>
//             </SearchBox>
//             <div className="option">
//               {/* <button type='button'>Submit</button> */}
//             </div>
//           </TopBar>
//           <LayoutBody>
//             <SideBar>
//               <RefinementListFilter id="tags" title="Tags" field="fieldTODO" size={10}/>
//             </SideBar>
//             <LayoutResults>
//               <ActionBar>
//                 <ActionBarRow>
//                   <GroupedSelectedFilters/>
//                 </ActionBarRow>
//                 <ActionBarRow>
//                   <HitsStats translations={{
//                     "hitstats.results_found": "{hitCount} results found"
//                   }}/>
//                 </ActionBarRow>
//               </ActionBar>
//             </LayoutResults>
//             <SideBar>
//             </SideBar>
//           </LayoutBody>

//         </Layout>
//       </SearchkitProvider>
//     );
//   }

// }
