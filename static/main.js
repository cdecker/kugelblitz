var stateColors = {
  STATE_NORMAL: "positive",
  STATE_ERR_BREAKDOWN: "negative",
  STATE_CLOSE_ONCHAIN_OUR_UNILATERAL: "negative",
  STATE_OPEN_WAITING_OURANCHOR: "warning",
  STATE_OPEN_WAITING_THEIRANCHOR: "warning",
  STATE_OPEN_WAITING_OURANCHOR_THEYCOMPLETED: "warning",
  STATE_ERR_INFORMATION_LEAK: "negative",
  STATE_NORMAL_COMMITTING: "positive"
};

function d3jsonrpc(url, method, args, cb) {
  params = {"method": method, "params": args, "jsonrpc": "2.0", "id": 0}
  d3.xhr('/rpc/').header("Content-Type", "application/json").header("Accept", "application/json")
    .post(JSON.stringify(params),
    function (error, data) {
      if (!error){
        data = JSON.parse(data.responseText);
        cb(data.error, data.result)
      } else {
        cb(error, null)
      }
    }
  );
}

function updatePeerTable() {
  d3jsonrpc('/rpc/', "LightningRpc.GetPeers", {}, function(error, data){
    if (error) {
      console.log(error);
    } else {
      var tbody = d3.select('#peersTbl > tbody');
      var rows = tbody.selectAll("tr").data(data.peers);
      rows.enter().append("tr");
      rows.exit().remove();
      rows.html(function(d) {
        return ("<td>" + 
                d.peerid + "</td><td>"+ d.connected+"</td><td>" + 
                d.state +"</td><td><button class='ui icon button open-connect-modal negative tiny disconnect-button' data-peerid='" + d.peerid + "'><i class='minus circle icon'></i> Disconnect</button></td></tr>");
      });
      rows.attr('class', function(d){ return stateColors[d.state]; });;
    }
  });
}

function updateInfo(){
  d3jsonrpc("/rpc/", "LightningRpc.GetInfo", {}, function(error, r){
    var row = d3.select("#nodeinfo tbody tr")
    var data = [r.id, r.version, r.port, r.testnet];      
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
  });

  d3jsonrpc("/rpc/", "Bitcoin.GetInfo", {}, function(error, r){
      var row = d3.select("#btcinfo tbody tr")
      var data = [r.version, r.blocks, r.connections, r.balance];
            row.selectAll("td").data(data).enter().append("td");
      row.selectAll("td").data(data).exit().remove();
      row.selectAll("td").text(function(e){return e;});
      if(r.error != null){
        $('#btcinfo').removeClass('green').addClass("red");
        data = [];
      }else{
        $('#btcinfo').removeClass("red").addClass("green");
      }

  });
}

function serializeFormData(form) {
  return form.serializeArray().reduce(function(obj, item) {
           obj[item.name] = item.value;
           return obj;
         }, {});
}

$('#peersTbl').on('click', '.disconnect-button', function(e) {
  var peerid = $(e.target).data('peerid');
  console.log("Disconnecting", peerid)
  d3jsonrpc('/rpc/', 'LightningRpc.Close', {"peerid": peerid}, function(error, data){
    console.log("Disconnected", data, error);
  });
});

var sendPaymentData = null;

$(document).ready(function(){
  window.setInterval(updatePeerTable, 10000);
  updatePeerTable()
  window.setInterval(updateInfo, 10000);
  updateInfo();


  $('.open-connect-modal').click(function(){
    $("#connect-dialog").modal("show");
  });
  $('#send-button').click(function(){
    $('#send-dialog').modal('show');
  });

  $('#send-dialog form').form({
    on: 'blur',
    fields: {
      destination: ['exactLength[66]'],
      paymenthash: ['exactLength[64]'],
      amount: ['integer[1..4000000000]']
    },
    onSuccess: function (e) {
      var form = $(e.target);
      sendPaymentData = {
        destination: form.find('input[name="destination"]').val(),
        amount: parseInt(form.find('input[name="amount"]').val()),
        paymenthash: form.find('input[name="paymenthash"]').val(),
        route: null
      };
      d3jsonrpc('/rpc/', 'LightningRpc.GetRoute', {
        amount: sendPaymentData.amount,
        destination: sendPaymentData.destination,
        risk: 1
      },function(error, data){
          if (error) {
            var errors = $(e.target).closest('.modal').find('.error').first();
            errors.show().append("<ul><li>Error computing route: " + error.message + "</li></ul>");
          } else {
            $(e.target).closest('.modal').modal('hide');
          }
          console.log(data);
        });
      return false;
    }
  });

  $('#receive-button').click(function(){
    $('#receive-dialog').modal('show');
  });

  /* Instead of closing the modal, look for a form in it and submit that instead */
  $('.modal').modal({
    onApprove : function(e) {
      $(e).closest(".modal").find("form").submit();
      return false;
    }
  });
  
  $('#connect-dialog form').form({
    on: 'blur',
    fields: {
      capacity: 'integer[100000..4000000]',
      host: ['empty'],
      port: 'integer[1..65535]'
    },
    onSuccess: function(e){
      var data = serializeFormData($('#connect-dialog form'));
      data.async = true;
      data.port = parseInt(data.port);
      data.capacity = parseInt(data.capacity);
      d3jsonrpc("/rpc/", "Node.ConnectPeer", data, function(error, data){});
      $(e.target).closest('.modal').modal('hide');
      return false;
      }
  });

}); /* EOF onload */