import * as React from 'react';
import * as lodash from 'lodash';
import Helmet from 'react-helmet';
import * as _ from 'lodash';

import {
  SearchkitManager,SearchkitProvider,
  SearchBox, RefinementListFilter, Pagination,
  HierarchicalMenuFilter, HitsStats, SortingSelector, NoHits,
  ResetFilters, RangeFilter, NumericRefinementListFilter,
  ViewSwitcherHits, ViewSwitcherToggle, DynamicRangeFilter,
  InputFilter, GroupedSelectedFilters,
  Layout, TopBar, LayoutBody, LayoutResults,
  ActionBar, ActionBarRow, SideBar, TagFilterConfig} from 'searchkit'
// import { Popup } from '../popup';
import {
  Button, Alert, Spinner, Modal, ModalHeader,
  ModalFooter, ModalBody, Form, FormField, FormInput,
  Checkbox, Radio, RadioGroup, Pill
} from 'elemental';

import axios from "axios";
import { WithContext as ReactTags } from 'react-tag-input';

import '../../../node_modules/font-awesome/css/font-awesome.min.css';

const host = "http://localhost:9200/telegram"
const searchkit = new SearchkitManager(host)

const HitsGridItem = (props)=> {
  return (
    <div>
      <h1>WIP</h1>
    </div>
  )
}

const HitsListItem = (props)=> {
  const {bemBlocks, result} = props
  // let photoUrl = "https://cdn5.telesco.pe/file/JMPBFOKtg7SARQveUVzY0sXSqk7pUF7Nc5sbHFNvviSWJy-LFjEigg9V7gC_xc-tW_XJnhOX7Rlkkeb3ZZ5nq1Nf_dMbOmTzxgtn44sF4LSlPU2pv5XfQxlfLSVAQOdaVziBdgHER7-SvNqpRMznVaAjZbq75X-PKS8nFFH2Vt30qiBnrQDEz6nXnunQVa5Jgzjizrh8lcCNvCQLIGArl66X10HOI2CvjKynhNenNcsOBW2BICJ1VYjtUDAoN5KZwePAekNhN8APpksDmfUvH-kCmzyzz1lUUyCMSRcYzs4xgKQSjC_7t6kTuT_O_3EnbChOkQq6h9opXo0PHyP4aw.jpg"
  // let photoUrl = "http://localhost:18080/images/" + result._source.tgid+ ".jpg"
  // TODO: get jpg from source
  let photoUrl = "http://localhost:18080/images/" + "knarfeh" + ".jpg"
  let tDotMe = "https://t.me/" + result._source.tgid
  var sectionStyle = {
    width: "100%",
    height: "122px",
  };
  const source = lodash.extend({}, result._source, result.highlight)
  var pills = [];
  source.tags.map(function(item) {
    pills.push(<Pill label={item.name} type="primary" key={item.name}/>)
  })

  return (
    <div className={bemBlocks.item().mix(bemBlocks.container("item"))} data-qa="hit">
      <div className={bemBlocks.item("poster")}>
        <img alt="presentation" data-qa="poster" style={sectionStyle} src={photoUrl}/>
      </div>
      <div className={bemBlocks.item("details")}>
        <div>
          <h2 className={bemBlocks.item("title")} dangerouslySetInnerHTML={{__html:source.title? source.title: source.tgid}}></h2>
        </div>
        <h3 className={bemBlocks.item("subtitle")}>{source.desc}</h3>
        {pills}
        <div>
          <a target="_blank" className="btn btn-outline-danger btn-sm" href={tDotMe}>
            <i className="fa fa-telegram fa-2x"></i>
          </a>
        </div>
      </div>
    </div>
  )
}

 export interface ReactTags {
  id: number;
  text: string
}

export default class Homepage extends React.Component<{}, { showPopup: boolean, tgID: string, type: string, tags: ReactTags, description: string, }> {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */
  constructor(props) {
    super(props);
    this.state = {
      showPopup: false,
      tgID: "",
      type: "",
      // tags: "", // TODO: define a interface
      // tags: [{ id: 1, text: "Thailand" }, { id: 2, text: "India" }],
      tags: [],
      description: "",
    };

    this.togglePopup = this.togglePopup.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);

    this.handleDelete = this.handleDelete.bind(this);
    this.handleAddition = this.handleAddition.bind(this);
    this.handleDrag = this.handleDrag.bind(this);
  }

  handleDelete(i) {
    let tags = this.state.tags;
    tags.splice(i, 1);
    this.setState({tags: tags});
  }

  handleAddition(tag) {
    let tags = this.state.tags;
    tags.push({
        id: tags.length + 1,
        text: tag
    });
    this.setState({tags: tags});
  }

  handleDrag(tag, currPos, newPos) {
    let tags = this.state.tags;

    // mutate array
    tags.splice(currPos, 1);
    tags.splice(newPos, 0, tag);

    // re-render
    this.setState({ tags: tags });
  }

  togglePopup() {
    this.setState({
      showPopup: !this.state.showPopup
    });
  }

  handleSubmit(event) {
    const self = this;
    event.preventDefault();
    const tags = _.map(this.state.tags, (item) => {
      return {
        name: item["text"],
        count: 1
      }
    })
    if(this.state.type !== "") {
      tags.push({
        name: this.state.type,
        count: 1
      })
    } else {
      tags.push({
        name: "unknown",
        count: 1
      })
    }
    const postData = {
      tgid: this.state.tgID,
      desc: this.state.description,
      tags: tags,
      type: this.state.type? this.state.type: "unknown",
    }
    axios({
      method: 'post',
      url: '/api/v1/tg',
      data: postData
    }).then(function (response) {
      // TODO, not user friendly
      setTimeout(() => {
        window.location.href = window.location.href
      }, 100);
    }) ;
    // TODO: popup hint, then refresh
    // add fefault type
    // handle not exist!!!!
  }

  handleInputChange(event) {
    if(typeof event === "string") {
      this.setState({
        ["type"]: event
      })
      return
    }
    const target = event.target;
    const value = target.value;
    const name = target.name;
    this.setState({
      [name]: value
    });
  }

  render() {
    const formInputStyle = {
      "paddingRight": 0
    }
    const { tags, } = this.state;
    return (
      <div className="app">
        <SearchkitProvider searchkit={searchkit}>
          <Layout>
            <TopBar>
              <div className="st-logo"><a href="/" className="st-link">Search Telegram</a></div>
              <SearchBox autofocus={true} searchOnChange={true} />
              <div className="option">
                <Modal isOpen={this.state.showPopup} onCancel={this.togglePopup} backdropClosesModal>
                  <ModalHeader text="Submit people/group/channel" showCloseButton onClose={this.togglePopup} />
                  <ModalBody>
                    <Form onSubmit={this.handleSubmit}>
                      <FormField label="Telegram ID">
                        <FormInput autoFocus={true} placeholder="Telegram ID" name="tgID" value={this.state.tgID} onChange={this.handleInputChange} style={formInputStyle} />
                      </FormField>
                      <FormField label="Type(Optional)">
                        <RadioGroup onChange={this.handleInputChange} value={this.state.type} name="type" options={[
                          {name: "people", label: "People", value: "people"},
                          {name: "group", label: "Group", value: "group"},
                          {name: "channel", label: "Channel", value: "channel"},
                          {name: "bot", label: "Bot", value: "bot"},
                        ]} inline={true}>
                        </RadioGroup>
                      </FormField>
                      <FormField label="Tags(Optional)">
                        {/* <FormInput placeholder="Tags" name="tags" onChange={this.handleInputChange} style={formInputStyle} /> */}
                        <div>
                          <ReactTags
                            tags={tags}
                            autofocus={false}
                            handleDelete={this.handleDelete}
                            handleAddition={this.handleAddition}
                            handleDrag={this.handleDrag} />
                        </div>
                      </FormField>
                      <FormField label="Description(Optionnal)">
                        <FormInput placeholder="Description" name="description" onChange={this.handleInputChange} style={formInputStyle} value={this.state.description} multiline />
                      </FormField>
                      <Button type="primary" onClick={this.togglePopup} submit>Submit</Button>
                    </Form>
                  </ModalBody>
                </Modal>
                <Button type="primary" onClick={this.togglePopup}>Submit New</Button>
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
                      // {key:"grid", title:"Grid", itemComponent: HitsGridItem},
                      {key:"list", title:"List", itemComponent: HitsListItem, defaultOption:true}
                    ]}
                    scrollTo="body"
                />
                <NoHits suggestionsField={"tgid"}/>
                <Pagination showNumbers={true}/>
              </LayoutResults>
            </LayoutBody>
          </Layout>
        </SearchkitProvider>
      </div>
    );
  }
}



