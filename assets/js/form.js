(function() {
  var forms = {
    signup: {
      selector: '#signup-form',
      endpoint: '/subscribe',
      errorResponseMessage: function(xhr) { return null; },
      serverErrorMessage: null,
      successMessage: null
    },
    question: {
      selector: '#question-form',
      endpoint: '/question',
      errorResponseMessage: function(xhr) { return xhr.response; },
      serverErrorMessage: 'Something went wrong. Please try again.',
      successMessage: 'Thank you for your question!'
    },
    rsvp: {
      selector: '#rsvp-form',
      endpoint: '/rsvp',
      errorResponseMessage: function (xhr) { return xhr.response; },
      serverErrorMessage: 'Something went wrong. Please try again.',
      successMessage: 'Thank you for your RSVP!'
    },
    admin: {
      selector: '#comm-form',
      endpoint: '/communicate',
      errorResponseMessage: function (xhr) { return xhr.response; },
      serverErrorMessage: 'Something went wrong. Please try again.',
      successMessage: 'Success!'
    }
  };

  var form;
  switch ($('form')[0].id) {
    case 'signup-form':
      form = forms.signup;
      break;
    case 'question-form':
      form = forms.question;
      break;
    case 'rsvp-form':
      form = forms.rsvp;
      break;
    case 'comm-form':
      form = forms.admin;
      break;
  }

  var $form = $(form.selector)[0];
  var $submit = $('input[type="submit"]', $form)[0];
  var $message = $('.message-container > span.message', $form)[0];

  if (!('addEventListener' in $form)) { return; }

  $message._show = function(type, text) {
    $message.classList.add(type, 'visible');
    $message.innerHTML = text;

    window.setTimeout(
      function() { $message._hide(); },
      5000
    );
  };

  $message._hide = function() {
    $message.classList.remove('visible');
  };

  $form.addEventListener('submit', function(event) {
    event.stopPropagation();
    event.preventDefault();

    $message._hide();
    $submit.disabled = true;
    $message.classList.add('disabled');

    var xhr = new XMLHttpRequest();
    xhr.open("POST", form.endpoint, true);
    xhr.onload = function(e) {
      $submit.disabled = false;

      if (xhr.status < 200 || xhr.status > 299) {
        $message._show('failure', form.errorResponseMessage(xhr));

      } else {
        $form.reset();
        $message._show('success', form.successMessage);
      }
    }

    xhr.onerror = function(e) {
      $message._show('failure', form.serverErrorMessage);
    }

    xhr.send(new FormData($form))
  });
})();
