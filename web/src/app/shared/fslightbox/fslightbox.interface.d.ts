interface Window {
  /**
   * After adding new <a> tags to detect new lightboxes or update existing ones use global refreshFsLightbox() function
   */
  refreshFsLightbox: () => void;

  /**
   * Every unique value of data-fslightbox will be treated as instance.
   * To access specific instance, use it's methods, attach events to it you need to use global fsLightboxInstances object.
   */
  fsLightboxInstances: {
    [key: string]: {
      /**
       * Close lightbox
       */
      close: () => void,

      props: {
        /**
         * Every time instance is opened (both show and initialize)
         */
        onOpen: () => void,

        /**
         * Every time instance is closed
         */
        onClose: () => void,
      }
    },
  };
}
