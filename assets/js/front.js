$(function(){
  // alertTimeout it's the default timeout
  // to show the alerts, 2s
  var alertTimeout = 2000

  // buildCustomEvent it's a function to prepend
  // all the custom events so they are separated
  // from the normal BB flow
  var buildCustomEvent = function(event) {
    return "notigator:"+event
  }

  // catchError it's a global function that will call showAlert
  // with the response error
  var catchError = function(collection, response, options) {
    showAlert(response.responseJSON.error)
  }


  var Notification = Backbone.Model.extend({
    attributes: {
      "id": "",
      "title": "",
      "url": "",
      "scope": "",
      "updated_at": "",
    },
  });
  var NotificationList = Backbone.Collection.extend({
    model: Notification,

    comparator: 'updated_at',

    parse: function(response, options) {
      return response.data;
    },

    initialize: function(models, options) {
      this.url = options.url
    },

    groupByScopes: function() {
      var gps = this.groupBy(function(e) {
        return e.get("scope");
      });

      var res = []

      _.each(gps, function(notifications, scope) {
        var ne =_.sortBy(notifications, function(event) {
          return event.get("updated_at");
        });
        ne = ne.reverse();
        res.push({
          scope: scope,
          notifications: ne,
          updated_at: ne[0].get("updated_at"),
        });
      });

      return _.sortBy(res, function(item) {
        return item.updated_at;
      }).reverse();
    },
  });

  var Source = Backbone.Model.extend({
    id: "canonical",
    attributes: {
      "canonical": "",
      "name": "",
      "active": false,
    },
    initialize: function(options) {
      this.notifications = new NotificationList(null, {
        url: this.url() + "/" + this.get("canonical") + "/notifications"
      })
    },
  });

  var Alert = Backbone.Model.extend({
    attributes: {
      "text": "",
    }
  })

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

  var NotificationsView = Backbone.View.extend({
    template: _.template($('#notifications-view-tmpl').html()),
    events: {
      "click .scope-title": "toggleNotifications",
    },
    render: function() {
      this.$el.html(this.template({ notifications: this.collection }));
      return this;
    },
    toggleNotifications: function(e) {
      this.$el.find("#notifications-"+e.currentTarget.id).toggle("fast")
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

  // showAlert prints an alert with the
  // given text
  var showAlert = function(text) {
    var a = new AlertView({model: new Alert({text: text})});
    $("#alert").html(a.render().el);
    setTimeout(function() {
      $("#alert").html("");
    }, alertTimeout);
  }

  var Sources = new SourceList();

  var AppView = Backbone.View.extend({
    el: $("#main"),

    initialize: function() {
      this.listenTo(Sources, 'sync', this.renderSources);
      this.listenTo(Sources, buildCustomEvent('change:active'), this.fetchNotifications);

      Sources.fetch({success: this.setDefaultSource.bind(this), error: catchError});
    },

    // setDefaultSource sets the sources collection
    // 'active' canonical to the actual model
    // withing it
    setDefaultSource: function(sources, response, options) {
      var that = this
      var src = sources.first()
      src.set('active', true);
      sources.active = src;
      sources.each(function(src) {
        that.listenTo(src.notifications, 'sync', that.renderNotifications);
        that.listenTo(src.notifications, 'reset', that.resetNotifications);
      })

      // I could just call the this.fetchNotifications but I feel that having
      // the actual event would help to understand the flow
      Sources.trigger(buildCustomEvent('change:active'));
    },

    renderNotifications: function() {
      var notifications = new NotificationsView({ collection: Sources.active.notifications });
      this.$("#list-notifications").html(notifications.render().el);
      return this;
    },

    renderSources: function() {
      Sources.each(function(source) {
        var s = new SourceView({model: source});
        this.$("#list-sources").append(s.render().el);
      });
      return this;
    },

    resetNotifications: function() {
      this.$("#list-notifications").html("");
    },

    fetchNotifications: function() {
      Sources.active.notifications.fetch({reset: true, error: catchError});
    },
  });

  var App = new AppView;
})
