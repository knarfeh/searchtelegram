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
  ActionBar, ActionBarRow, SideBar, TagFilterConfig, Hits} from 'searchkit'
import {
  Button, Alert, Spinner, Modal, ModalHeader,
  ModalFooter, ModalBody, Form, FormField, FormInput,
  Checkbox, Radio, RadioGroup, Pill, EmailInputGroup
} from 'elemental';
import axios from "axios";
import { WithContext as ReactTags } from 'react-tag-input';
import { ToastContainer, toast } from 'react-toastify';
import '../../../node_modules/font-awesome/css/font-awesome.min.css';

const host = "http://localhost:18080/"   // TODO: configurable
const searchkit = new SearchkitManager(host)

const HitsListItem = (props)=> {
  const {bemBlocks, result} = props
  let tDotMe = "https://t.me/" + result._source.tgid
  // let photoUrl = result._source.imgsrc
  let photoUrl = "http://localhost:18080" + result._source.imgsrc
  // let photoUrl = "https://s3.amazonaws.com/searchtelegram/media/images/telegram.jpg"
  var sectionStyle = {
    width: "122px",
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

export interface ReactTagsItem {
  id: number;
  text: string
}

export default class Homepage extends React.Component<{}, { showPopup: boolean, tgID: string, type: string, tags: ReactTags[], description: string, validateTgIDMessage: any}> {
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
      validateTgIDMessage: "",
    };
    this.togglePopup = this.togglePopup.bind(this);
    this.handleInputChange = this.handleInputChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);

    this.handleDelete = this.handleDelete.bind(this);
    this.handleAddition = this.handleAddition.bind(this);
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

  togglePopup() {
    this.setState({
      showPopup: !this.state.showPopup
    });
  }

  handleSubmit(event) {
    const self = this;
    event.preventDefault();
    if(this.state.tgID === "") {
      this.setState({
        validateTgIDMessage: (
          <div className="form-validation is-invalid">
            Telegram ID is required
          </div>
        )
      })
      return
    }

    this.togglePopup()
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
    }
    const postData = {
      tgid: this.state.tgID,
      desc: this.state.description,
      tags: tags,
      type: this.state.type,
    }
    axios({
      method: 'post',
      url: '/api/v1/tg',
      data: postData
    }).then(function (response) {
      toast("If everything goes well, you will be able to search for it after a while.");
      setTimeout(() => {
        searchkit.reloadSearch();
      }, 5000);
    }) ;
    // unit test, e2e test
    // limit api frequency, after online
    // add detail page
    // add update page
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
    if(name === "tgID" && value !== "") {
      this.setState({
        validateTgIDMessage: ""
      })
    }
    this.setState({
      [name]: value
    });
  }

  render() {
    const formInputStyle = {
      "paddingRight": 0
    };
    const { tags, } = this.state;
    return (
      <div className="app">
        <ToastContainer />
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
                      <FormField label="Telegram ID, please make sure it exist">
                        <FormInput autoFocus={true} placeholder="Telegram ID" name="tgID" value={this.state.tgID} onChange={this.handleInputChange} style={formInputStyle} />
                        {this.state.validateTgIDMessage}
                      </FormField>
                      <FormField label="Type, optional">
                        <RadioGroup onChange={this.handleInputChange} value={this.state.type} name="type" options={[
                          {name: "people", label: "People", value: "people"},
                          {name: "group", label: "Group", value: "group"},
                          {name: "channel", label: "Channel", value: "channel"},
                          {name: "bot", label: "Bot", value: "bot"},
                        ]} inline={true}>
                        </RadioGroup>
                      </FormField>
                      <FormField label="Tags, optional">
                        <div>
                          <ReactTags
                            tags={tags}
                            autofocus={false}
                            handleDelete={this.handleDelete}
                            handleAddition={this.handleAddition} />
                        </div>
                      </FormField>
                      <FormField label="Description, optionnal">
                        <FormInput placeholder="Description" name="description" onChange={this.handleInputChange} style={formInputStyle} value={this.state.description} multiline />
                      </FormField>
                      <Button type="primary" submit>Submit</Button>
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
                      "hitstats.results_found":"{hitCount} results were found"
                    }}/>
                    <ViewSwitcherToggle/>
                  </ActionBarRow>
                  <ActionBarRow>
                    <GroupedSelectedFilters/>
                    <ResetFilters/>
                  </ActionBarRow>
                </ActionBar>
                <Hits
                  hitsPerPage={10}
                  highlightFields={["tgid", "title"]}
                  sourceFilter={["tgid", "title", "desc", "type", "tags", "imgsrc"]}
                  mod="sk-hits-list"
                  itemComponent={HitsListItem}
                />
                <NoHits suggestionsField={"tgid"}/>
                <Pagination showNumbers={true}/>
              </LayoutResults>
            </LayoutBody>
          </Layout>
        </SearchkitProvider>
        <div className="container">
          <div id="footer" className="row">
            <div className="col-xs-12 col-md-5">
              <div> © 2018 <a href="https://github.com/knarfeh" target="_blank"> knarfeh </a> |
                <a href="#" target="_blank"> Source </a> |
                <a href="#"> Telegram Bot </a> |
                <a href="#"> Email </a>
              </div>
              <br/>
            </div>
            <div className="col-md-1 text-center">
            </div>
            <div className="clearfix visible-sm visible-xs"><br/></div>
            <div className="col-xs-12 col-md-5 text-right">
              Donate BTC: <a className="pointer" data-toggle="modal" data-target="#donate_btc">3CMCRgEm8HVz3DrWaCCid3vAANE42jcEv9</a><br/>
              Donate ETH: <a className="pointer" data-toggle="modal" data-target="#donate_eth">0x0074709077B8AE5a245E4ED161C971Dc4c3C8E2B</a><br/>
              <a className="pointer" data-toggle="modal" data-target="#donate_ada">Donate ADA</a><br/>
            </div>
          </div>
        </div>
      </div>
    );
  }
}



