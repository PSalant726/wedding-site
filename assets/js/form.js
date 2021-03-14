(function() {
  var forms = {
    signup: {
      selector: '#signup-form',
      endpoint: '/subscribe',
      errorResponseMessage: function(xhr) { return null; },
      serverErrorMessage: null,
      successMessage: null,
      availableText: "Subscribe",
      waitingText: "Subscribing...",
    },
    question: {
      selector: '#question-form',
      endpoint: '/question',
      errorResponseMessage: function(xhr) { return xhr.response; },
      serverErrorMessage: 'Something went wrong. Please try again.',
      successMessage: 'Thank you for your question!',
      availableText: "Ask",
      waitingText: "Asking...",
    },
    rsvp: {
      selector: '#rsvp-form',
      endpoint: '/rsvp',
      errorResponseMessage: function (xhr) { return xhr.response; },
      serverErrorMessage: 'Something went wrong. Please try again.',
      successMessage: 'Success! Thanks for your reply!',
      availableText: "Submit RSVP",
      waitingText: "Submitting...",
    },
    admin: {
      selector: '#comm-form',
      endpoint: '/communicate',
      errorResponseMessage: function (xhr) { return xhr.response; },
      serverErrorMessage: 'Something went wrong. Please try again.',
      successMessage: 'Success!',
      availableText: "Send Communication",
      waitingText: "Sending...",
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
      30000
    );
  };

  $message._hide = function() {
    $message.classList.remove('disabled', 'failure', 'visible');
    $message.innerHTML = " ";
  };

  $form.addEventListener('submit', function(event) {
    event.stopPropagation();
    event.preventDefault();

    $message._hide();
    $submit.disabled = true;
    $submit.value = form.waitingText
    $submit.classList.add('disabled')
    $message.classList.add('disabled');

    var xhr = new XMLHttpRequest();
    xhr.open("POST", form.endpoint, true);
    xhr.onload = function(e) {
      $submit.disabled = false;
      $submit.classList.remove('disabled')
      $submit.value = form.availableText

      if (xhr.status < 200 || xhr.status > 299) {
        $message._show('failure', form.errorResponseMessage(xhr));

      } else {
        $form.reset();
        $message._show('success', form.successMessage);
      }
    }

    xhr.onerror = function(e) {
      $submit.disabled = false;
      $submit.classList.remove('disabled')
      $submit.value = form.availableText
      $message._show('failure', form.serverErrorMessage);
    }

    xhr.send(new FormData($form))
  });
})();
