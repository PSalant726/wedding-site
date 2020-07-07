(function () {
  var $form = document.querySelectorAll('#signup-form')[0];
  var $list = document.querySelectorAll('#signup-form ul')[0];
  var $submit = document.querySelectorAll('#signup-form input[type="submit"]')[0];
  var $message = document.createElement('li');

  if (!('addEventListener' in $form)) { return; }

  $message.appendChild(document.createElement("span"));
  $($message, 'span').addClass("message");
  $list.appendChild($message);

  $message._show = function (type) {
    $message.classList.add(type);
    $message.classList.add('visible');

    window.setTimeout(
      function () { $message._hide(); },
      3000
    );
  };

  $message._hide = function () {
    $message.classList.remove('visible');
  };

  $form.addEventListener('submit', function (event) {
    event.stopPropagation();
    event.preventDefault();

    $message._hide();
    $submit.disabled = true;

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/subscribe", true);
    xhr.onload = function (e) {
      $form.reset();
      $submit.disabled = false;
      $message._show('success');
    }

    xhr.onerror = function (e) {
      $message._show('failure');
    }

    xhr.send(new FormData($form))
  });
})();
