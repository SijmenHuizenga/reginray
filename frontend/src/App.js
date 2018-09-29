import React, {Component} from 'react';
import './App.css';
import {TopBar} from './components/TimeSelector'
import {LogTable} from './components/LogTable'
import InfiniteScroll from 'react-infinite-scroller';
import PropTypes from 'prop-types';
import history from './history'

const pageSize = 25;

class Root extends Component {

  constructor(props) {
    super(props);
    this.state = {
      begin: unixTimestampSeconds() - (6 * 24 * 60 * 60),  //5 dag
      end: unixTimestampSeconds(),
    }
  }

  render() {
    return <App onChangeTimerange={this.onChangeTimerange.bind(this)} logsStart={this.state.begin} logsEnd={this.state.end}/>
  }

  onChangeTimerange(begin, end) {
    history.push({search: "?begin="+begin+"&end="+end});
    this.setState({
      begin: begin,
      end: end,
    })
  }

}

class App extends Component {

  static propTypes = {
    logsStart: PropTypes.number,
    logsEnd: PropTypes.number,
  };

  constructor(props) {
    super(props);

    this.state = {
      logs : [],
      nextLogsPage : 0,
      moreLogsAvailable : true,
      logsLoading: false,
      logsPerTimevalue : [],
    };
  }

  render() {
    return (

      <div className="App">
        <header className="header">
          <h1 className="title">Aspicio</h1>
          <div className="timeselector">
            <TopBar
              logsPerTimevalue={this.state.logsPerTimevalue}
              startTime={this.props.logsStart}
              finishTime={this.props.logsEnd}
              interval={this.getInterval()}
              timeRangeChanged={this.timeRangeChanged.bind(this)}
            />
          </div>
        </header>
        <div className="body">
          <InfiniteScroll
            pageStart={0}
            loadMore={this.loadMore.bind(this)}
            hasMore={this.state.moreLogsAvailable}
            loader={<div className="loader" key={0}>Loading ...</div>}
          >
            <LogTable logs={this.state.logs}/>
          </InfiniteScroll>
        </div>
      </div>
    );
  }

  timeRangeChanged(begin, end) {
    if (end - begin < 120) {
      end = begin + 120;
    }

    this.props.onChangeTimerange(begin, end);
  }

  loadMore() {
    if(this.state.logsLoading)
      return;
    this.setState({
      logsLoading: true,
    }, () => {
      this.fetchLogs()
    })
  }

  componentDidMount() {
    this.fetchMetrics();
  }

  componentDidUpdate(prevProps, a, b) {
    if (this.props.logsEnd !== prevProps.logsEnd || this.props.logsStart !== prevProps.logsStart) {
      this.setState({
        logs : [],
        nextLogsPage : 0,
        moreLogsAvailable : true,
      }, () => {
        this.fetchMetrics()
      });

    }
  }

  fetchLogs() {
    fetch("http://localhost:8080/logs?timeperiod=" + this.props.logsStart + "-" + this.props.logsEnd + "&pagesize="+pageSize+"&page=" + this.state.nextLogsPage)
      .then(response => response.json())
      .then(responseJson => {
        if (responseJson.hasOwnProperty("error")) {
          this.setState({logsLoading: false});
          alert(responseJson.error);
        } else {
          this.setState({
            logs : [...this.state.logs, ...responseJson],
            logsLoading : false,
            nextLogsPage: this.state.nextLogsPage + 1,
            moreLogsAvailable : responseJson.length === pageSize
          })
        }
      }).catch(alert);
  }

  fetchMetrics() {
    fetch("http://localhost:8080/stats?interval=" + this.getInterval() + "&timeperiod=" + this.props.logsStart + "-" + this.props.logsEnd)
      .then(response => response.json())
      .then(responseJson => {
        if (responseJson.hasOwnProperty("error")) {
          alert(responseJson.error);
        } else {
          this.setState({logsPerTimevalue : responseJson})
        }
      }).catch(alert);
  }

  getInterval(){
    return Math.round((this.props.logsEnd - this.props.logsStart) / 120);
  }
}

function unixTimestampSeconds() {
  return Math.round((new Date()).getTime() / 1000)
}

export default (Root);
