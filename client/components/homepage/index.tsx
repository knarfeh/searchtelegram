import * as React from 'react';
import * as lodash from 'lodash';
import Helmet from 'react-helmet';

import { SearchkitManager,SearchkitProvider,
  SearchBox, RefinementListFilter, Pagination,
  HierarchicalMenuFilter, HitsStats, SortingSelector, NoHits,
  ResetFilters, RangeFilter, NumericRefinementListFilter,
  ViewSwitcherHits, ViewSwitcherToggle, DynamicRangeFilter,
  InputFilter, GroupedSelectedFilters,
  Layout, TopBar, LayoutBody, LayoutResults,
  ActionBar, ActionBarRow, SideBar, TagFilterConfig} from 'searchkit'

import '../../../node_modules/font-awesome/css/font-awesome.min.css';

const host = "http://localhost:9200/telegram"
const searchkit = new SearchkitManager(host)

const MovieHitsGridItem = (props)=> {
  const {bemBlocks, result} = props
  let url = "https://cdn5.telesco.pe/file/qyqzjBHidDTCg5MWywQn5hHdpZkZvDRZnD9578Up785eEO2AXtzkPOFgHd0AK5TFgoNwaaJdv8lxQwfF-GsxrjaUpS_kdIOQtCVLD7QEllGg3d-PZ466DWzUHI8dgEyeXJgpCtOKMd0OnA4Ziyv1-ZulKozHv9t9OUgx0GKbZ7gac3-xTYx9S9y5k90XDh4N4dJmALLQaoLgBUDbDENeKAPOsSk0wnVdWHkG879wd2MRnQouYdnXldv2lIdXcQOdYj9J66uuRSx_X27O2Go3QjTYeP7pMtUz7BUCyos3YOOQOqB_xl_y7I4w84C3MHjv360Om5uFBT9mtyJL8iyR7A.jpg"
  const source = lodash.extend({}, result._source, result.highlight)
  return (
    <div className={bemBlocks.item().mix(bemBlocks.container("item"))} data-qa="hit">
      <a href={url} target="_blank">
        <img data-qa="poster" alt="presentation" className={bemBlocks.item("poster")} src={url} width="170" height="240"/>
        <div data-qa="title" className={bemBlocks.item("title")} dangerouslySetInnerHTML={{__html:source.name}}></div>
      </a>
    </div>
  )
}

const MovieHitsListItem = (props)=> {
  const {bemBlocks, result} = props
  // let photoUrl = "https://cdn5.telesco.pe/file/JMPBFOKtg7SARQveUVzY0sXSqk7pUF7Nc5sbHFNvviSWJy-LFjEigg9V7gC_xc-tW_XJnhOX7Rlkkeb3ZZ5nq1Nf_dMbOmTzxgtn44sF4LSlPU2pv5XfQxlfLSVAQOdaVziBdgHER7-SvNqpRMznVaAjZbq75X-PKS8nFFH2Vt30qiBnrQDEz6nXnunQVa5Jgzjizrh8lcCNvCQLIGArl66X10HOI2CvjKynhNenNcsOBW2BICJ1VYjtUDAoN5KZwePAekNhN8APpksDmfUvH-kCmzyzz1lUUyCMSRcYzs4xgKQSjC_7t6kTuT_O_3EnbChOkQq6h9opXo0PHyP4aw.jpg"
  // let photoUrl = "http://localhost:18080/images/" + result._source.tgid+ ".jpg"
  let photoUrl = "http://localhost:18080/images/" + "knarfeh" + ".jpg"
  let tDotMe = "https://t.me/" + result._source.tgid
  var sectionStyle = {
    width: "100%",
    height: "122px",
  };
  const source = lodash.extend({}, result._source, result.highlight)
  return (
    <div className={bemBlocks.item().mix(bemBlocks.container("item"))} data-qa="hit">
      <div className={bemBlocks.item("poster")}>
        <img alt="presentation" data-qa="poster" style={sectionStyle} src={photoUrl}/>
      </div>
      <div className={bemBlocks.item("details")}>
        <h2 className={bemBlocks.item("title")} dangerouslySetInnerHTML={{__html:source.title? source.title: source.tgid}}></h2>
        <h3 className={bemBlocks.item("subtitle")}>@{source.tgid}</h3>
        <h3 className={bemBlocks.item("subtitle")}>Tags: "test"</h3>
        <div className="TODO">
          <a target="_blank" className="btn btn-outline-danger btn-sm" href={tDotMe}>
            <i className="fa fa-telegram fa-3x"></i>
          </a>
        </div>
      </div>
    </div>
  )
}

export class Popup extends React.Component<{ text:string, closePopup:any }, { }> {
  render() {
    return (
      <div className='popup'>
        <div className='popup_inner'>
          <h1>{this.props.text}</h1>
        <button onClick={this.props.closePopup}>close me</button>
        </div>
      </div>
    );
  }
}

export default class Homepage extends React.Component<{}, { showPopup: any }> {
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
    this.state = {
      showPopup: false
    };
  }
  togglePopup() {
    this.setState({
      showPopup: !this.state.showPopup
    });
    console.log('showpopup???', this.state.showPopup)
  }

  render() {
    return (

      <div className="app">
        {this.state.showPopup ? <Popup text='Close Me' closePopup={this.togglePopup.bind(this)}/> : null}
      <SearchkitProvider searchkit={searchkit}>
        <Layout>
          <TopBar>
            <div className="st-logo">Search Telegram</div>
            <SearchBox autofocus={true} searchOnChange={true} />
            <div className="option">
              <button className="btn btn-sm btn-default" onClick={this.togglePopup.bind(this)}>show popup</button>
            </div>
          </TopBar>

        <LayoutBody>

          <SideBar>
            <RefinementListFilter id="tags" title="tags" field="tags.name.keyword" operator="OR" size={10}/>
          </SideBar>
          <LayoutResults>
            <ActionBar>

              <ActionBarRow>
                <HitsStats translations={{
                  "hitstats.results_found":"{hitCount} results found"
                }}/>
                <ViewSwitcherToggle/>
              </ActionBarRow>

              <ActionBarRow>
                <GroupedSelectedFilters/>
                <ResetFilters/>
              </ActionBarRow>

            </ActionBar>
            <ViewSwitcherHits
                hitsPerPage={12} highlightFields={["tgid", "title", "info"]}
                sourceFilter={["tgid", "title", "info", "desc", "type", "tags"]}
                hitComponents={[
                  {key:"grid", title:"Grid", itemComponent:MovieHitsGridItem},
                  {key:"list", title:"List", itemComponent:MovieHitsListItem, defaultOption:true}
                ]}
                scrollTo="body"
            />
            <NoHits suggestionsField={"name"}/>
            <Pagination showNumbers={true}/>
          </LayoutResults>

          </LayoutBody>
        </Layout>
      </SearchkitProvider>
      </div>
    );
  }
}



