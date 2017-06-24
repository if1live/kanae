function sync(url, type) {
  fetch(url, {
    method: 'post'
  }).then(function (resp) {
    alert('sync ' + type + ' success');

  }).catch(function (err) {
    console.log(err);
    alert('sync ' + type + ' failed');
  });
}
document.querySelector('.btn-sync-trade').onclick = function () {
  sync('/sync/trade', 'trade');
}

document.querySelector('.btn-sync-lending').onclick = function () {
  sync('/sync/lending', 'lending');
}

document.querySelector('.btn-sync-balance').onclick = function () {
  sync('/sync/balance', 'balance');
}
