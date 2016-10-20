var stateColors = {
  STATE_NORMAL: "positive",
  STATE_ERR_BREAKDOWN: "negative",
  STATE_CLOSE_ONCHAIN_OUR_UNILATERAL: "negative",
  STATE_OPEN_WAITING_OURANCHOR: "warning",
  STATE_OPEN_WAITING_OURANCHOR_THEYCOMPLETED: "warning",
  STATE_ERR_INFORMATION_LEAK: "negative",
  STATE_NORMAL_COMMITTING: "positive"
};

function updatePeerTable() {
  params = {"method": "LightningRpc.GetPeers", "params": [], "jsonrpc": "2.0", "id": 0}

  d3.xhr('/rpc/').header("Content-Type", "application/json")
    .post(JSON.stringify(params),
    function (error, data) {
      data = JSON.parse(data.responseText);

    var columns = {
      'name': {
        'text': function(obj, c){return obj[c]}
      },
      'connected': {
        'text': function(obj, c){return obj[c]}
      },
      'state': {
        'text': function(obj, c){return obj[c]}
      }
    };
    
    var tbody = d3.select('#peersTbl > tbody');
    var rows = tbody.selectAll("tr").data(data.result.peers);
    rows.enter().append("tr");
    rows.exit().remove();
    rows.html(function(d) {
      return ("<td>" + 
      d.peerid + "</td><td>"+ d.connected+"</td><td>" + 
      d.state  + "</td></tr>");
    });
    rows.attr('class', function(d){ return stateColors[d.state]; });
  });
}

function updateInfo(){
  params = {"method": "LightningRpc.GetInfo", "params": [], "jsonrpc": "2.0", "id": 0}
  d3.xhr('/rpc/').header("Content-Type", "application/json")
    .post(JSON.stringify(params),
    function(error, r){
      r = JSON.parse(r.responseText)
      var row = d3.select("#nodeinfo tbody tr")
      var data = [r.result.id, r.result.version, r.result.port, r.result.testnet];      
      if(r.error != null){
        $('#nodeinfo').removeClass('green').addClass("red");
        $("#connection-lost").show()
        data = [];
      }else{
        $('#nodeinfo').removeClass("red").addClass("green");
        $("#connection-lost").hide()
      }

      row.selectAll("td").data(data).enter().append("td");
      row.selectAll("td").data(data).exit().remove();
      row.selectAll("td").text(function(e){return e;});
    }
  );
  params = {"method": "Bitcoin.GetInfo", "params": [], "jsonrpc": "2.0", "id": 0}
  d3.xhr('/rpc/').header("Content-Type", "application/json")
    .post(JSON.stringify(params),
    function(error, r){
      r = JSON.parse(r.responseText)
      console.log(r)
  });
}

function serializeFormData(form) {
  return form.serializeArray().reduce(function(obj, item) {
           obj[item.name] = item.value;
           return obj;
         }, {});
}

$(document).ready(function(){
  window.setInterval(updatePeerTable, 1000);
  updatePeerTable()
  window.setInterval(updateInfo, 5000);
  updateInfo();


  $('.open-connect-modal').click(function(){
    $("#connect-dialog").modal("show");
  });

  $('#connect-dialog').modal({
    onApprove: function(e){
      if($("#connect-dialog .error").length > 0)
        return false;

      var data = serializeFormData($('#connect-dialog form'));
      $.post(
        "/rpc/echo",{
          data: JSON.stringify({
            id: 1,
            method: "connect",
            params: data
          })
        },
        function(e){}
      );
      return true;
    }
  });

  $('#connect-dialog form').form({
    on: 'blur',
    fields: {
      capacity: {
        identifier: 'capacity',
        rules: [{
          type   : 'integer[100000..4000000]',
          prompt : ''
        }]
      },
      address: {
        identifier: 'address',
        rules: [
          {
            type   : 'empty',
            prompt : 'Please enter a node address'
          }
        ]
      },
      port: {
        identifier: 'port',
        rules: [
          {
            type   : 'integer[1..65535]',
            prompt : 'This is not a valid port.'
          }
        ]
      }
    }
  });
});
