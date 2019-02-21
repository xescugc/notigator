$(function(){
  // defaultSourceCanonical it's the default
  // source that will always be used
  var defaultSourceCanonical = 'github'

  // alertTimeout it's the default timeout
  // to show the alerts, 2s
  var alertTimeout = 2000

  // buildCustomEvent it's a function to prepend
  // all the custom events so they are separated
  // from the normal BB flow
  var buildCustomEvent = function(event) {
    return "notigator:"+event
  }

  var Event = Backbone.Model.extend({
    attributes: {
      "id": "",
      "title": "",
      "url": "",
      "scope": "",
      "updated_at": "",
    },
  });
  var Source = Backbone.Model.extend({
    id: 'canonical',
    attributes: {
      "canonical": "",
      "name": "",
      "active": false,
    },
  });
  var Alert = Backbone.Model.extend({
    attributes: {
      "text": "",
    }
  })

  var EventList = Backbone.Collection.extend({
    model: Event,
    url: function() {
      var active = Sources.active ? Sources.active.get('canonical') : defaultSourceCanonical
      return '/api/sources/'+active+'/notifications';
    },

    comparator: 'updated_at',

    parse: function(response, options) {
      return response.data;
    },

    groupByScopes: function() {
      var gps = this.groupBy(function(e) {
        return e.get("scope");
      });

      var res = []

      _.each(gps, function(events, scope) {
        var ne =_.sortBy(events, function(event) {
          return event.get("updated_at");
        });
        ne = ne.reverse();
        res.push({
          scope: scope,
          events: ne,
          updated_at: ne[0].get("updated_at"),
        });
      });

      return _.sortBy(res, function(item) {
        return item.updated_at;
      }).reverse();

    },
  });

  var SourceList = Backbone.Collection.extend({
    model: Source,
    url: '/api/sources',

    initialize: function() {
      // Initialize the attribute
      this.active = undefined;
    },

    parse: function(response, options) {
      return response.data;
    },
  });

  var Events = new EventList();
  var Sources = new SourceList();

  var EventsView = Backbone.View.extend({
    template: _.template($('#events-view-tmpl').html()),
    render: function() {
      this.$el.html(this.template({ events: this.collection }));
      return this;
    },
  })

  var SourceView = Backbone.View.extend({
    tagName: 'a',
    template: _.template($('#source-view-tmpl').html()),
    className: 'list-group-item list-group-item-action',
    events: {
      "click": "changeActive",
    },
    initialize: function() {
      this.listenTo(this.model, "change", this.render);
    },
    attributes: function(){
      return {
        'href': '#' + this.model.get('canonical'),
      };
    },
    changeActive: function(e) {
      if (!this.model.get("active")) {
        this.model.set('active',true);
        Sources.active.set('active', false);
        Sources.active = this.model;
        Sources.trigger(buildCustomEvent('change:active'));
      };
    },
    render: function() {
      if (this.model.get('active')) {
        $(this.el).addClass('active').html(this.template(this.model.toJSON()));
      } else {
        $(this.el).removeClass('active').html(this.template(this.model.toJSON()));
      }
      return this;
    },
  });

  var AlertView = Backbone.View.extend({
    className: 'alert alert-danger',
    template: _.template($('#alert-view-tmpl').html()),
    model: Alert,
    render: function() {
      this.$el.html(this.template(this.model.toJSON()));
      return this;
    }
  })

  // setDefaultSource sets the sources collection
  // 'active' canonical to the actual model
  // withing it
  var setDefaultSource = function(sources, response, options) {
    var src = sources.findWhere({canonical: defaultSourceCanonical});
    src.set('active', true);
    sources.active = src;
  }

  // showAlert prints an alert with the
  // given text
  var showAlert = function(text) {
    var a = new AlertView({model: new Alert({text: text})});
    $("#alert").html(a.render().el);
    setTimeout(function() {
      $("#alert").html("");
    }, alertTimeout);
  }

  // catchError it's a global function that will call showAlert
  // with the response error
  var catchError = function(collection, response, options) {
    showAlert(response.responseJSON.error)
  }

  var AppView = Backbone.View.extend({
    el: $("#main"),

    initialize: function() {
      this.listenTo(Events, 'sync', this.renderEvents);
      this.listenTo(Events, 'reset', this.resetEvents);
      this.listenTo(Sources, 'sync', this.renderSources);
      this.listenTo(Sources, buildCustomEvent('change:active'), this.fetchEvents);

      Sources.fetch({success: setDefaultSource, error: catchError});
      Events.fetch({error: catchError});
    },

    renderEvents: function() {
      var events = new EventsView({ collection: Events });
      this.$("#list-events").html(events.render().el);
      return this;
    },

    renderSources: function() {
      Sources.each(function(source) {
        var s = new SourceView({model: source});
        this.$("#list-sources").append(s.render().el);
      });
      return this;
    },

    resetEvents: function() {
      this.$("#list-events").html("");
    },

    fetchEvents: function() {
      Events.fetch({reset: true, error: catchError});
    },
  });

  var App = new AppView;
})
