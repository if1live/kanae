function sync(el, url, type) {
  el.disabled = true;

  fetch(url, {
    method: 'post'
  }).then(function (resp) {
    el.disabled = false;
    alert('sync ' + type + ' success');

  }).catch(function (err) {
    el.disabled = false;
    console.log(err);
    alert('sync ' + type + ' failed');
  });
}

document.querySelector('.btn-sync-trade').onclick = function (evt) {
  sync(this, '/sync/trade', 'trade');
}

document.querySelector('.btn-sync-lending').onclick = function (evt) {
  sync(this, '/sync/lending', 'lending');
}

document.querySelector('.btn-sync-balance').onclick = function (evt) {
  sync(this, '/sync/balance', 'balance');
}
