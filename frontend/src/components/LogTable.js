import React, {Component} from 'react';
import './Table.css'
import Moment from 'moment'
import {default as AnsiUp} from 'ansi_up';
import {ContextMenu, ContextMenuProvider, Item} from "react-contexify";

const ansi_up = new AnsiUp();

//takes {} object and returns [[key, value], [key, value]]
function recursiveEntries(obj) {
  return Object.entries(obj).flatMap(entrie => {
    //If a value is an object replace it with prefixed child objects.
    if (typeof entrie[1] === 'object') {
      return recursiveEntries(entrie[1])
        .map(newentrie => [entrie[0] + "." + newentrie[0], newentrie[1]])
    }
    return [entrie]
  });
}

//takes [[key, value], [key, value]] and returns {key: value, key: value}
function objectify(inn) {
  let out = {};
  for (let i = 0; i < inn.length; i++) {
    out[inn[i][0]] = inn[i][1]
  }
  return out
}


class LogText extends Component {
  render() {
    let fields = this.props.selectedFields;

    if (fields.length === 0) {
      return <span>{ansi_up.ansi_to_html(this.props.log.Message)}</span>
    }

    let tagList = recursiveEntries(this.props.log.Fields)
      .filter((el) => fields.includes(el[0]));

    return tagList.flatMap(field =>
        [<span className="tag" key={field[0]}>{field[0]}:</span>, " " + field[1] + " "]
      )
  }
}

export class LogTable extends Component {

  constructor() {
    super();
    this.state = {
      selectedFields : []
    }
  }

  render() {

    return <div className="logs">
      <ContextMenuProvider id="tablecontextmenu">
        <table>
          <tbody>
          {this.props.logs.map(log =>
            <tr key={log.Id}>
              <td>{this.formatTimestamp(log.TimestampSeconds)} <span className="tog">[{log.Container.Name}]</span></td>
              <td><LogText log={log} selectedFields={this.state.selectedFields}/></td>
            </tr>
          )}
          </tbody>
        </table>
      </ContextMenuProvider>
      {this.renderMenu()}
    </div>
  }

  renderMenu() {
    let uniqueFieldList = this.uniqueFieldsList();

    let cl = this.onClick.bind(this);
    let f = this.state.selectedFields;
    return <ContextMenu id='tablecontextmenu' theme={"dark"}>
      <Item onClick={cl} data={{action: "set", value: []}} disabled={f.length === 0}>Show no tags</Item>
      <Item onClick={cl} data={{action: "set", value: uniqueFieldList}} disabled={uniqueFieldList.length == f.length}>Show all tags</Item>
      {uniqueFieldList.map(
        (key) => {
          if(f.includes(key)) {
            return <Item onClick={cl} data={{action: "remove", value: key}} key={key}>- {key}</Item>
          }else {
            return <Item onClick={cl} data={{action: "add", value: key}} key={key}>+ {key}</Item>
          }
        }
      )}

    </ContextMenu>
  }

  onClick({event, ref, data }) {
    switch (data.action) {
      case "set":
        this.setState({
          selectedFields: data.value
        });
        break;
      case "add":
        this.setState({
          selectedFields: [...this.state.selectedFields, data.value]
        });
        break;
      case "remove":
        this.setState({
          selectedFields: this.state.selectedFields.filter((x) => data.value !== x)
        });
        break;
    }
  }

  formatTimestamp(unixtimestampseconds) {
    let t = Moment(unixtimestampseconds * 1000);
    return t.format("YYYY-MM-DD HH:mm:ss");
  }

  uniqueFieldsList(){
    let out = [];
    for(let i = 0; i < this.props.logs.length; i++)
      Object.keys(this.props.logs[i].Fields).forEach((key) => {
        if(out.indexOf(key) < 0)
          out.push(key)
      })
    return out;
  }

}