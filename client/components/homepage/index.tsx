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

const host = "https://searchtelegram.com"   // TODO: configurable
const searchkit = new SearchkitManager(host)

const HitsListItem = (props)=> {
  const {bemBlocks, result} = props
  let tDotMe = "https://t.me/" + result._source.tgid
  // let photoUrl = result._source.imgsrc
  let photoUrl = "https://searchtelegram.com" + result._source.imgsrc
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

export default class Homepage extends React.Component
  <{}, {
    showSubmitPopup: boolean,
    showBTCPopup: boolean,
    showETHPopup: boolean,
    showADAPopup: boolean,
    tgID: string,
    type: string,
    tags: ReactTags[],
    description: string,
    validateTgIDMessage: any
  }> {
  /*eslint-disable */
  static onEnter({store, nextState, replaceState, callback}) {
    // Load here any data.
    callback(); // this call is important, don't forget it
  }
  /*eslint-enable */
  constructor(props) {
    super(props);
    this.state = {
      showSubmitPopup: false,
      showBTCPopup: false,
      showETHPopup: false,
      showADAPopup: false,
      tgID: "",
      type: "",
      tags: [],
      description: "",
      validateTgIDMessage: "",
    };
    this.toggleSubmitPopup = this.toggleSubmitPopup.bind(this);
    this.toggleBTCPopup = this.toggleBTCPopup.bind(this);
    this.toggleETHPopup = this.toggleETHPopup.bind(this);
    this.toggleADAPopup = this.toggleADAPopup.bind(this);

    this.handleDonateClick = this.handleDonateClick.bind(this);
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

  toggleSubmitPopup() {
    this.setState({
      showSubmitPopup: !this.state.showSubmitPopup
    });
  }

  toggleBTCPopup() {
    this.setState({
      showBTCPopup: !this.state.showBTCPopup
    })
  }

  toggleETHPopup() {
    this.setState({
      showETHPopup: !this.state.showETHPopup
    })
  }

  toggleADAPopup() {
    this.setState({
      showADAPopup: !this.state.showADAPopup
    })
  }

  handleDonateClick(event) {
    const type = event.target.getAttribute('data-target')
    switch(type) {
      case 'btc': {
        this.setState({showBTCPopup: !this.state.showBTCPopup});
        break;
      }
      case 'eth': {
        this.setState({showETHPopup: !this.state.showETHPopup});
        break;
      }
      case 'ada': {
        this.setState({showADAPopup: !this.state.showADAPopup});
        break;
      }
    };
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

    this.toggleSubmitPopup()
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
      }, 7000);
    }).catch(function (error) {
      if(error.response.status == 503) {
        toast.error("Maybe you submitted too often, please try again later. " +
        " Perhaps you can support the server through donations")
      }
    }) ;
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
                <Modal isOpen={this.state.showSubmitPopup} onCancel={this.toggleSubmitPopup} backdropClosesModal>
                  <ModalHeader text="Submit people/group/channel" showCloseButton onClose={this.toggleSubmitPopup} />
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
                <Button type="primary" onClick={this.toggleSubmitPopup}>New Submit</Button>
              </div>
            </TopBar>
            <LayoutBody>
              <SideBar>
                <RefinementListFilter id="tags" title="tags" field="tags.name.keyword" operator="OR" size={10}/>
              </SideBar>
              <LayoutResults>
                <ActionBar>
                  <ActionBarRow>
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
                <a href="mailto:knarfeh@outlook.com"> Email </a> |
                <a href="https://github.com/knarfeh/searchtelegram" target="_blank"> Source </a> |
                <a href="https://t.me/SearchTelegramPublic" target="_blank"> Telegram Group </a>
              </div>
              <br/>
            </div>
            <div className="col-md-1 text-center">
            </div>
            <div className="clearfix visible-sm visible-xs"><br/></div>
            <div className="col-xs-12 col-md-5 text-right">
            Donate <a className="pointer" data-toggle="modal" data-target="btc" onClick={this.handleDonateClick}>BTC </a> |
              <a className="pointer" data-toggle="modal" data-target="eth" onClick={this.handleDonateClick}> ETH </a> |
              <a className="pointer" data-toggle="modal" data-target="ada" onClick={this.handleDonateClick}> ADA</a>
            </div>
          </div>
          <Modal isOpen={this.state.showBTCPopup} onCancel={this.toggleBTCPopup} backdropClosesModal>
            <ModalHeader text="Donate BTC" showCloseButton onClose={this.toggleBTCPopup} />
            <ModalBody>
              <div className="text-center">
                <p>1Aa8ZXPbzoyHGp9SmnWjSaNq56py3jCj96</p>
                <br/>
                <img src="/images/searchtelegrambitcoinqrcode.png" alt="Donate Bitcoin" />
              </div>
            </ModalBody>
          </Modal>
          <Modal isOpen={this.state.showETHPopup} onCancel={this.toggleETHPopup} backdropClosesModal>
            <ModalHeader text="Donate ETH" showCloseButton onClose={this.toggleETHPopup} />
            <ModalBody>
              <div className="text-center">
                <p>0x3A149665Fb7fe1b44892D50eCA0bd2BdcD21C85D</p>
                <br/>
                <img src="/images/searchtelegramethqrcode.png" alt="Donate Ethereum" />
              </div>
            </ModalBody>
          </Modal>
          <Modal isOpen={this.state.showADAPopup} onCancel={this.toggleADAPopup} backdropClosesModal width="large">
            <ModalHeader text="Donate ADA" showCloseButton onClose={this.toggleADAPopup} />
            <ModalBody>
              <div className="text-center">
                <p>DdzFFzCqrhssJHqka9HYyaXXYPdasFQEaSPRULDxgkgBN7jkkDmueznLaP3cMNuouLE4WHs5cLrHtn8oCn3FAJQ7fa2eytjzxS4CmH6K</p>
                <br/>
                <img src="/images/searchtelegramdonateqrcode.png" alt="Donate Ada coin" />
              </div>
            </ModalBody>
          </Modal>
        </div>
      </div>
    );
  }
}



