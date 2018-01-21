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

var updateInterval = 1000;

var sendPaymentData = null;

function d3jsonrpc(url, method, args, cb) {
  params = {"method": method, "params": args, "jsonrpc": "2.0", "id": 0}
  d3.xhr('/rpc/').header("Content-Type", "application/json").header("Accept", "application/json")
    .post(JSON.stringify(params),
    function (error, data) {
      if (!error){
        data = JSON.parse(data.responseText);
        cb(null, data.error, data.result)
      } else {
        cb(error, null, null)
      }
    }
  );
}

function updateHistory() {
    d3jsonrpc('/rpc/', "Node.GetHistory", {}, function(terror, error, data){
	console.log(data);
    var tbody = d3.select('#historyTbl > tbody');
    var rows = tbody.selectAll("tr");
    if (error) {
      rows.remove();
    } else if(terror) {
      transportFailure(terror);
    } else {

      rows = rows.data(data);
      rows.enter().append("tr");
	rows.exit().remove();
	rows.html(function(d) {
	    console.log(d);
	    return (
		"<td>" + d.Destination + "</td>" +
		    "<td>" + d.Msatoshi + "</td>" +
		    "<td>" + d.Status + "</td>"
	    );
      });
    }
    });
}

function updatePeerTable() {
  d3jsonrpc('/rpc/', "Lightning.GetPeers", {}, function(terror, error, data){
    var tbody = d3.select('#peersTbl > tbody');
    var rows = tbody.selectAll("tr");
    if (error) {
      rows.remove();
    } else if(terror) {
      transportFailure(terror);
    } else {
      rows = rows.data(data.peers);
      rows.enter().append("tr");
      rows.exit().remove();
      rows.html(function(d) {
        return ("<td>" +
                d.peerid.substring(0,15) + "...</td><td>"+ d.connected+"</td><td>" +
                d.state +"</td><td><button class='ui icon button open-connect-modal negative tiny disconnect-button' data-peerid='" + d.peerid + "'><i class='minus circle icon'></i> Disconnect</button></td></tr>");
      });
      rows.attr('class', function(d){ return stateColors[d.state]; });;
    }
  });
}

function transportFailure(terror) {
      setAllState('red', "Connection to <em>kugelblitz</em> lost, can't check other daemons.")
      d3.select("#nodeinfolist").selectAll(".item").remove();
}

function updateLightningInfo() {
  d3jsonrpc("/rpc/", "Lightning.GetInfo", {}, function(terror, error, r){
    var headers = ["Node ID", "Version", "Port", "Testnet"]
    var items = d3.select("#nodeinfolist").selectAll(".item");
    if (terror){
      transportFailure(terror);
      items.remove();
    } else if(error) {
      setLightningState('red', "Connection to <em>lightningd</em> lost.")
      items.remove();
    } else {
      var data = [r.id, r.version, r.port, r.testnet];
      items.data(data).enter().append("div").classed('item', true);
      items.data(data).exit().remove();
      items = d3.select("#nodeinfolist").selectAll(".item");
      items.html(function(e, t){return "<div class='header'>"+headers[t]+"</div>" + e});
      setLightningState('green', "Lightningd is up and running.");
    }
  });
}

function updateBitcoinInfo() {
  d3jsonrpc("/rpc/", "BitcoinRpc.GetInfo", {}, function(terror, error, r){
    if (error){
      console.log("Error retrieving bitcoind info", error);
      setBitcoinState('red', "Could not contact <em>bitcoind</em>, maybe we just need to wait?")
      return;
    }else{
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
      d3jsonrpc("/rpc/", "Node.GetFundingAddr", {}, function(terror, error, data){
        console.log(error, data.addr);
        $("#btc-fund-addr").html(data.addr);
      });
      $("#btc-no-funds-error").show();
      $('#btcinfo').removeClass('green').addClass("red").addClass("attached");
      $("#btc-error").html(error).hide();
      setBitcoinState('yellow', "Your bitcoin node does not have any funds available. We can't create channels without funds.")
    }else{
      $('#btcinfo').removeClass("red").removeClass("attached").addClass("green")
      $("#btc-no-funds-error").hide();
      $("#btc-error").html(error).hide();
      info.bitcoin = r;
      setBitcoinState('green', "Your bitcoin node is up and running.")
    }
    }
  });
}

function updateKugelblitzInfo() {
  console.log("Updating kugelblitz")
  d3jsonrpc("/rpc/", "Node.GetInfo", {}, function(terror, error, r){
    console.log(terror, error, r);
    if (error || terror){
      console.log("Error retrieving kugelblitz info", error);
      setKugelblitzState('red', "Could not retrieve information from Kugelblitz: " + error);
      return;
    } else {
      setKugelblitzState('green', "Kugelblitz is up and running.")
    }
  });
}

function updateInfo(){
  updateLightningInfo();
  updateBitcoinInfo();
  updateKugelblitzInfo();
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
  d3jsonrpc('/rpc/', 'Lightning.Close', {"peerid": peerid}, function(terror, error, data){
    console.log("Disconnected", data, error);
    updateInfo();
  });
  $(e.target).addClass('disabled');
});

$(document).ready(function(){
  window.setInterval(updatePeerTable, updateInterval);
  updatePeerTable()
  window.setInterval(updateInfo, updateInterval);
  updateInfo();
//  window.setInterval(updateInfo, updateHistory);
    updateHistory();

  $('.open-connect-modal').click(function(){
    $("#connect-dialog").modal("show");
  });
  $('#send-button').click(function(){
    $('#route-info').hide();
    $('#send-dialog').modal('show');
  });

  $('#send-dialog form').form({
    on: 'blur',
    onSuccess: function (e) {
      var form = $(e.target);
      //window.sendPaymentData = {
        //destination: form.find('input[name="destination"]').val(),
        //amount: parseInt(form.find('input[name="amount"]').val()),
        //paymenthash: form.find('input[name="paymenthash"]').val(),
        route: window.sendPaymentData.route
      //};
    console.log(window.sendPaymentData)
    $('#send-dimmer').addClass('active');
    d3jsonrpc('/rpc/', 'Lightning.SendPayment', {
      route: sendPaymentData.route,
	paymenthash: sendPaymentData.paymenthash
    }, function(terror, error, data){
          $('#send-dimmer').removeClass('active');
         if (!error){
             $('#send-dialog').modal('hide');
	     updateHistory();
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
      destination: form.find('input[name="destination"]').val()
    };
    d3jsonrpc('/rpc/', 'Lightning.GetPaymentRequestInfo', {
      destination: window.sendPaymentData.destination
    },function(terror, error, data){

        if (error) {
          var errors = $(e.target).closest('.modal').find('.error').first();
          //errors.empty().append("<ul><li>Error computing route: " + error.message + "</li></ul>").show();
        } else {
	    window.sendPaymentData.paymenthash = data.paymenthash;
	    window.sendPaymentData.amount = data.amount;
            showRoute(data.route);
	    $('#route-destination').text("");
	    $('#route-payment-hash').text(data.paymenthash);
	    $('#route-amount').text(data.amount);
	    $(window).trigger('resize');
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
	port: 'integer[1..65535]',
	nodeid: ['empty']
    },
    onSuccess: function(e){
      var data = serializeFormData($('#connect-dialog form'));
      data.async = true;
      data.port = parseInt(data.port);
      data.capacity = parseInt(data.capacity);
      d3jsonrpc("/rpc/", "Node.ConnectPeer", data, function(terror, error, data){});
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
    {id: "This node (source)", delay: 0, msatoshi: 0}
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

var color2state = {
    red: 'negative',
    yellow: 'warning',
    green: 'positive'
}
function setDaemonState(element, color, status) {
  element.find('span.fa-stack').removeClass('yellow red green').addClass(color);
  element.find('.popup .message')
  .html(status).removeClass('warning positive negative').addClass(color2state[color]);
}
