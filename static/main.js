var stateColors = {
  STATE_NORMAL: "positive",
  STATE_ERR_BREAKDOWN: "negative",
  STATE_CLOSE_ONCHAIN_OUR_UNILATERAL: "negative",
  STATE_OPEN_WAITING_OURANCHOR: "warning",
  STATE_OPEN_WAITING_THEIRANCHOR: "warning",
  STATE_OPEN_WAIT_ANCHORDEPTH_AND_THEIRCOMPLETE: "warning",
  STATE_OPEN_WAITING_OURANCHOR_THEYCOMPLETED: "warning",
  STATE_ERR_INFORMATION_LEAK: "negative",
  STATE_NORMAL_COMMITTING: "positive"
};

var info = {
  lightning: {},
  bitcoin: {}
};

var sendPaymentData = null;

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
    if (error){
      setAllState('red', "Connection to <em>kugelblitz</em> lost, can't check other daemons.")
      d3.select("#nodeinfolist").selectAll(".item").remove();
      //setLightningState('red', "Error retrieving <em>lightningd</em> info");
      console.log("Error retrieving lightningd info", error);
      return;
    }
    setLightningState('green', "Lightningd is up and running.");

    var data = [r.id, r.version, r.port, r.testnet];
    var headers = ["Node ID", "Version", "Port", "Testnet"]
    
    var items = d3.select("#nodeinfolist").selectAll(".item");
    items.data(data).enter().append("div").classed('item', true);
    items.data(data).exit().remove();
    items.html(function(e, t){return "<div class='header'>"+headers[t]+"</div>" + e});

    var row = d3.select("#nodeinfo tbody tr")
    if(error != null){
      $('#nodeinfo').removeClass('green').addClass("red");
      $("#connection-lost").show()
      data = [];
    }else{
      $('#nodeinfo').removeClass("red").addClass("green");
      $("#connection-lost").hide()
      info.lightning = r;
    }
    
    row.selectAll("td").data(data).enter().append("td");
    row.selectAll("td").data(data).exit().remove();
    row.selectAll("td").text(function(e){return e;});
  });

  d3jsonrpc("/rpc/", "BitcoinRpc.GetInfo", {}, function(error, r){
    if (error){
      console.log("Error retrieving bitcoind info", error);
      return;
    }
    var row = d3.select("#btcinfo tbody tr")
    var data = [r.version, r.blocks, r.connections, r.balance];
    row.selectAll("td").data(data).enter().append("td");
    row.selectAll("td").data(data).exit().remove();
    row.selectAll("td").text(function(e){return e;});
    if(error != null){
      $('#btcinfo').removeClass('green').addClass("red").addClass("attached");
      $("#btc-no-funds-error").hide();
      $("#btc-error").html(error).show();
      data = [];
    }else if(r.balance == 0){
      d3jsonrpc("/rpc/", "Node.GetFundingAddr", {}, function(error, data){
        console.log(error, data.addr);
        $("#btc-fund-addr").html(data.addr);
      });
      $("#btc-no-funds-error").show();
      $('#btcinfo').removeClass('green').addClass("red").addClass("attached");
      $("#btc-error").html(error).hide();
    }else{
      $('#btcinfo').removeClass("red").removeClass("attached").addClass("green") 
      $("#btc-no-funds-error").hide();
      $("#btc-error").html(error).hide();
      info.bitcoin = r;
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
    updateInfo();
  });
  $(e.target).addClass('disabled');
});

$(document).ready(function(){
  window.setInterval(updatePeerTable, 2000);
  updatePeerTable()
  window.setInterval(updateInfo, 2500);
  updateInfo();

  $('.open-connect-modal').click(function(){
    $("#connect-dialog").modal("show");
  });
  $('#send-button').click(function(){
    $('#route-info').hide();
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
      window.sendPaymentData = {
        destination: form.find('input[name="destination"]').val(),
        amount: parseInt(form.find('input[name="amount"]').val()),
        paymenthash: form.find('input[name="paymenthash"]').val(),
        route: window.sendPaymentData.route
      };
    console.log(window.sendPaymentData)
    $('#send-dimmer').addClass('active');
    d3jsonrpc('/rpc/', 'LightningRpc.SendPayment', {
      route: sendPaymentData.route,
      paymenthash: sendPaymentData.paymenthash
    }, function(error, data){
          $('#send-dimmer').removeClass('active');
         if (!error){
           $('#send-dialog').modal('hide');
         } else {
            var errors = $(e.target).closest('.modal').find('.error').first();
            errors.empty().append("<ul><li>Error sending payment: " + error.message + "</li></ul>").show();
           console.log(error);
         }
       });
      return false;
    }
  });

  $('#send-form input').on('blur', function(e) {
    var form = $(e.target).closest('form');
    window.sendPaymentData = {
      destination: form.find('input[name="destination"]').val(),
      amount: parseInt(form.find('input[name="amount"]').val()),
      paymenthash: form.find('input[name="paymenthash"]').val()
    };
    d3jsonrpc('/rpc/', 'LightningRpc.GetRoute', {
      amount: window.sendPaymentData.amount,
      destination: window.sendPaymentData.destination,
      risk: 1
    },function(error, data){
        
        if (error) {
          var errors = $(e.target).closest('.modal').find('.error').first();
          errors.empty().append("<ul><li>Error computing route: " + error.message + "</li></ul>").show();
        } else {
          window.sendPaymentData.route = data.route
          showRoute(data.route);
        }
      });
    return false;
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

  $('.fa-5x').popup({
    inline     : true,
    hoverable  : true,
    position   : 'bottom center',
    delay: {
      show: 100,
      hide: 250
    },
    hoverable: true
  });
  installHandlers();
}); /* EOF onload */

function showRoute(route) {
  window.sendPaymentData.route = route;
  var hops = [
    {id: info.lightning.id + " (source)", delay: 0, msatoshi: 0}
  ]
  $.each(route, function(_, e){
    hops.push(e);
    //e.delay *= 6;
  });
  var body = $('#route-info tbody').empty();

  $.each(hops, function(i, e){
    var pos = ""
    if (i == 0)
      pos = "-first";
    else if (i == hops.length - 1)
      pos = '-last';
            
    body.append('<tr><td class="route-segment"><svg viewBox="0 0 43 43"><use xlink:href="#route-segment'+pos+'"></svg></td><td>' + e.id + '</td><td> ' +e.msatoshi+ ' </td></tr>');
  });

  //$('#route-info').transition('slide down');
  $('#route-info').show();
}

function installHandlers() {
}

function setAllState(color, status) {
  setBitcoinState(color, status);
  setKugelblitzState(color, status);
  setLightningState(color, status);
}

function setBitcoinState(color, status) {
  return setDaemonState($('#btc-state'), color, status);
}

function setLightningState(color, status) {
  return setDaemonState($('#lightning-state'), color, status);
}

function setKugelblitzState(color, status) {
  return setDaemonState($('#kb-state'), color, status);
}

function setDaemonState(element, color, status) {
  element.find('span.fa-stack').removeClass('yellow red green').addClass(color);
  element.find('.popup .message').html(status);
}