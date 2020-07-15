(function () {
  var $form = document.querySelectorAll('#question-form')[0];
  var $list = document.querySelectorAll('#question-form > div.row')[0];
  var $submit = document.querySelectorAll('#question-form input[type="submit"]')[0];
  var $message = document.createElement('div');

  if (!('addEventListener' in $form)) { return; }

  $message.appendChild(document.createElement("span"));
  $($message, 'span').addClass("message");
  $list.appendChild($message);

  $message._show = function (type, text) {
    $message.innerHTML = text;
    $message.classList.add(type, 'visible', 'col-12');

    window.setTimeout(
      function () { $message._hide(); },
      5000
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
    $message.classList.add('disabled');

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/question", true);
    xhr.onload = function (e) {
      $form.reset();
      $submit.disabled = false;

      if (xhr.status != 200) {
        $message._show('failure', xhr.response);
      } else {
        $message._show('success', 'Thank you for your question!');
      }
    }

    xhr.onerror = function (e) {
      $message._show('failure', 'Something went wrong. Please try again.');
    }

    xhr.send(new FormData($form))
  });
})();
