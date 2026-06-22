import type { Editor, BlockProperties } from 'grapesjs';
import { RequiredPluginOptions } from '.';

export default (editor: Editor, opts: RequiredPluginOptions) => {
  const addBlock = (id: string, def: BlockProperties) => {
    opts.blocks.indexOf(id)! >= 0 && editor.Blocks.add(id, {
      select: true,
      category: 'Basic',
      ...def,
      ...opts.block(id),
    });
  }

  // Basic Row Block
  addBlock('row', {
    label: 'Row',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <rect x="2" y="10" width="20" height="4" fill="none" stroke="currentColor" stroke-width="1.5"/>
    <rect x="2" y="6" width="20" height="2" fill="currentColor" opacity="0.3"/>
    <rect x="2" y="16" width="20" height="2" fill="currentColor" opacity="0.3"/>
  </svg>`,
    content: `
    <div data-gjs-type="row" class="flex flex-wrap w-full" 
         data-gjs-editable="false" 
         data-gjs-droppable=".column" 
         data-gjs-name="Row">
      <div data-gjs-type="column" class="flex-1 min-w-[100px] min-h-[50px]" 
           data-gjs-droppable="true"
           data-gjs-name="Column">
      </div>
    </div>
  `,
  });

  // Single Column Block
  addBlock('column', {
    label: 'Column',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <rect x="8" y="2" width="8" height="20" fill="none" stroke="currentColor" stroke-width="1.5"/>
    <rect x="4" y="2" width="2" height="20" fill="currentColor" opacity="0.3"/>
    <rect x="18" y="2" width="2" height="20" fill="currentColor" opacity="0.3"/>
  </svg>`,
    content: `
    <div data-gjs-type="column" class="flex-1 min-w-[100px] min-h-[120px]" 
         data-gjs-droppable="true" 
         data-gjs-name="Column">
    </div>
  `,
  });

  // Row with 2 Columns
  addBlock('row-2-columns', {
    label: '2 Columns',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <rect x="2" y="4" width="9" height="16" fill="none" stroke="currentColor" stroke-width="1.5"/>
    <rect x="13" y="4" width="9" height="16" fill="none" stroke="currentColor" stroke-width="1.5"/>
  </svg>`,
    content: `
    <div data-gjs-type="row" class="flex flex-wrap w-full" 
         data-gjs-editable="false" 
         data-gjs-droppable=".column" 
         data-gjs-name="2 Column Row">
      <div data-gjs-type="column" class="flex-1 min-w-[200px] min-h-[100px]" 
           data-gjs-droppable="true" 
           data-gjs-name="Column 1">
      </div>
      <div data-gjs-type="column" class="flex-1 min-w-[200px] min-h-[100px]" 
           data-gjs-droppable="true" 
           data-gjs-name="Column 2">
      </div>
    </div>
  `,
  });

  // Row with 3 Columns
  addBlock('row-3-columns', {
    label: '3 Columns',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <rect x="2" y="4" width="5" height="16" fill="none" stroke="currentColor" stroke-width="1.5"/>
    <rect x="9" y="4" width="6" height="16" fill="none" stroke="currentColor" stroke-width="1.5"/>
    <rect x="17" y="4" width="5" height="16" fill="none" stroke="currentColor" stroke-width="1.5"/>
  </svg>`,
    content: `
    <div data-gjs-type="row" class="flex flex-wrap w-full" 
         data-gjs-editable="false" 
         data-gjs-droppable=".column" 
         data-gjs-name="3 Column Row">
      <div data-gjs-type="column" class="flex-1 min-w-[150px] min-h-[100px]" 
           data-gjs-droppable="true" 
           data-gjs-name="Column 1">
      </div>
      <div data-gjs-type="column" class="flex-1 min-w-[150px] min-h-[100px]" 
           data-gjs-droppable="true" 
           data-gjs-name="Column 2">
      </div>
      <div data-gjs-type="column" class="flex-1 min-w-[150px] min-h-[100px]" 
           data-gjs-droppable="true" 
           data-gjs-name="Column 3">
      </div>
    </div>
  `,
  });

  // RTL-Aware Row (for Arabic layouts)
  addBlock('row-rtl', {
    label: 'RTL Row',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <rect x="2" y="10" width="20" height="4" fill="none" stroke="currentColor" stroke-width="1.5"/>
    <path d="M20 8 L16 10 L20 12" stroke="currentColor" stroke-width="1.5" fill="none"/>
  </svg>`,
    content: `
    <div data-gjs-type="row" class="flex flex-wrap flex-row-reverse w-full gap-2" 
         dir="rtl"
         data-gjs-editable="false" 
         data-gjs-droppable=".column" 
         data-gjs-name="RTL Row">
      <div data-gjs-type="column" class="flex-1 min-w-[200px] min-h-[100px]" 
           data-gjs-droppable="true" 
           data-gjs-name="العمود الأول">
      </div>
      <div data-gjs-type="column" class="flex-1 min-w-[200px] min-h-[100px]" 
           data-gjs-droppable="true" 
           data-gjs-name="العمود الثاني">
      </div>
    </div>
  `,
  });

  // Responsive Container Block
  addBlock('container', {
    label: 'Container',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <rect x="3" y="5" width="18" height="14" fill="none" stroke="currentColor" stroke-width="1.5" rx="2"/>
    <rect x="1" y="3" width="22" height="2" fill="currentColor" opacity="0.3" rx="1"/>
    <rect x="1" y="19" width="22" height="2" fill="currentColor" opacity="0.3" rx="1"/>
  </svg>`,
    content: `
    <div data-gjs-type="container" class="container mx-auto px-4 max-w-6xl min-h-[200px] border border-dashed border-purple-200" 
         data-gjs-droppable="true" 
         data-gjs-name="Container">
      <div class="text-gray-400 text-sm text-center mt-16" data-gjs-type="text">
        Responsive Container - Drop content here
      </div>
    </div>
  `,
  });

  // Section Block
  addBlock('section', {
    label: 'Section',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <rect x="2" y="6" width="20" height="12" fill="none" stroke="currentColor" stroke-width="1.5" rx="1"/>
    <line x1="2" y1="10" x2="22" y2="10" stroke="currentColor" stroke-width="0.5" opacity="0.5"/>
    <line x1="2" y1="14" x2="22" y2="14" stroke="currentColor" stroke-width="0.5" opacity="0.5"/>
  </svg>`,
    content: `
    <section data-gjs-type="section" class="w-full py-12 px-4 min-h-[250px] bg-gray-50 border border-dashed border-gray-300" 
             data-gjs-droppable="true" 
             data-gjs-name="Section">
      <div class="container mx-auto">
        <div class="text-gray-400 text-center mt-8" data-gjs-type="text">
          <div class="text-lg font-medium mb-2">Section</div>
          <div class="text-sm">Drop your content here</div>
        </div>
      </div>
    </section>
  `,
  });

  // Text Block (Arabic/English ready)
  addBlock('text', {
    label: 'Text',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <path d="M4 5h16M4 12h16M4 19h10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
  </svg>`,
    content: `
    <div data-gjs-type="text" class="p-4 text-gray-800 leading-relaxed" 
         data-gjs-editable="true" 
         data-gjs-name="Text">
      <p>Enter your text here. يمكنك كتابة النص باللغة العربية أيضاً.</p>
    </div>
  `,
  });

  // Heading Block
  addBlock('heading', {
    label: 'Heading',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
    <path d="M4 6h16M4 10h12M4 14h16M4 18h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
    <rect x="2" y="4" width="20" height="4" fill="currentColor" opacity="0.1" rx="1"/>
  </svg>`,
    content: `
    <div data-gjs-type="text" class="p-4" 
         data-gjs-editable="true" 
         data-gjs-name="Heading">
      <h2 class="text-3xl font-bold text-gray-900 mb-4">Your Heading Here</h2>
      <h2 class="text-3xl font-bold text-gray-900 mb-4" dir="rtl">عنوانك هنا</h2>
    </div>
  `,
  });
  addBlock('row', {
    label: 'Row',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
  <rect x="2" y="10" width="20" height="8" fill="none" stroke="currentColor" stroke-width="1"/>
</svg>

`,
    content: `
        <div data-gjs-type="row" class="row" style="display:flex;flex-wrap:wrap;">
    <div data-gjs-type="column" class="column" style="flex:1;min-width:100px;min-height:50px;">
    </div>
    </div>
  `,
  });

  addBlock('column', {
    label: 'Column',
    category: 'Basic',
    media: `<svg viewBox="0 0 24 24">
        <path fill="currentColor" d="M2 20h20V4H2v16Zm-1 0V4a1 1 0 0 1 1-1h20a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1Z"></path>
      </svg>`,
    content: `
    <div data-gjs-type="column" class="column" style="flex:1;min-width:100px;min-height:50px;">
    </div>
  `,
  });
  // Blocks (note the data-gjs-type attributes)
  addBlock('column2', {
    label: 'Row & 2 Columns',
    category: 'Basic',
    media: `<svg viewBox="0 0 23 24">
        <path fill="currentColor" d="M2 20h8V4H2v16Zm-1 0V4a1 1 0 0 1 1-1h8a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1ZM13 20h8V4h-8v16Zm-1 0V4a1 1 0 0 1 1-1h8a1 1 0 0 1 1 1v16a1 1 0 0 1-1 1h-8a1 1 0 0 1-1-1Z"></path>
      </svg>`,
    content: `
    <div data-gjs-type="row" class="row" style="display:flex;flex-wrap:wrap;">
      <div  data-gjs-type="column" before="drop here" class="column" style="flex:1;min-width:100px;min-height:50px;">
      </div>
      <div data-gjs-type="column" before="drop here" class="column" style="flex:1;min-width:100px;min-height:50px;">
      </div>
    </div>
  `,
  });


  // HERO SECTIONS
  addBlock('ticket-hero', {
    label: 'Ticket Hero',
    category: 'Ticket Landing',
    media: '<i class="fa fa-ticket" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
    <section class="relative overflow-hidden bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900">
      <div class="absolute inset-0 opacity-10">
        <div class="absolute inset-0" style="background-image: url('data:image/svg+xml,<svg width="60" height="60" viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg"><g fill="none" fill-rule="evenodd"><g fill="%23ffffff" fill-opacity="0.1"><circle cx="30" cy="30" r="2"/></g></g></svg>'); background-size: 60px 60px;"></div>
      </div>
      <div class="relative mx-auto max-w-7xl px-6 py-16 lg:py-24">
        <div class="grid lg:grid-cols-2 gap-12 items-center">
          <div class="text-center lg:text-left">
            <div class="inline-flex items-center px-4 py-2 rounded-full bg-yellow-400 text-black text-sm font-semibold mb-6">
              <i class="fa fa-calendar mr-2"></i>
              Limited Time Event
            </div>
            <h1 class="text-4xl lg:text-6xl font-black tracking-tight text-white leading-tight">
              Get Your 
              <span class="text-transparent bg-clip-text bg-gradient-to-r from-yellow-400 to-orange-500">Tickets</span>
              <br>Before They're Gone
            </h1>
            <p class="mt-6 text-xl text-gray-300 leading-relaxed">
              Don't miss out on the event of the year. Secure your spot now with instant confirmation and mobile tickets.
            </p>
            <div class="mt-8 grid grid-cols-1 sm:grid-cols-3 gap-4 text-center lg:text-left">
              <div class="flex items-center justify-center lg:justify-start text-gray-300">
                <i class="fa fa-calendar-alt text-yellow-400 mr-3"></i>
                <div>
                  <div class="font-semibold text-white">March 15</div>
                  <div class="text-sm">2024</div>
                </div>
              </div>
              <div class="flex items-center justify-center lg:justify-start text-gray-300">
                <i class="fa fa-clock text-yellow-400 mr-3"></i>
                <div>
                  <div class="font-semibold text-white">8:00 PM</div>
                  <div class="text-sm">Doors at 7:30</div>
                </div>
              </div>
              <div class="flex items-center justify-center lg:justify-start text-gray-300">
                <i class="fa fa-map-marker-alt text-yellow-400 mr-3"></i>
                <div>
                  <div class="font-semibold text-white">City Arena</div>
                  <div class="text-sm">Downtown</div>
                </div>
              </div>
            </div>
            <div class="mt-10 flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4">
              <a href="#tickets" class="w-full sm:w-auto px-8 py-4 rounded-xl bg-gradient-to-r from-yellow-400 to-orange-500 text-black font-bold text-lg hover:from-yellow-300 hover:to-orange-400 transition-all duration-300 transform hover:scale-105 shadow-xl">
                <i class="fa fa-ticket-alt mr-2"></i>Buy Tickets Now
              </a>
              <a href="#info" class="w-full sm:w-auto px-8 py-4 rounded-xl border-2 border-white/30 text-white font-semibold hover:bg-white/10 transition-all duration-300">
                Event Details
              </a>
            </div>
          </div>
          <div class="flex justify-center lg:justify-end">
            <div class="relative">
              <div class="bg-white rounded-2xl shadow-2xl border border-gray-200 p-8 transform rotate-3 hover:rotate-0 transition-transform duration-500">
                <div class="flex justify-between items-start mb-6">
                  <div>
                    <div class="text-2xl font-bold text-gray-900">ADMIT ONE</div>
                    <div class="text-gray-600">General Admission</div>
                  </div>
                  <div class="text-right">
                    <div class="text-sm text-gray-500">Ticket #</div>
                    <div class="font-mono font-bold">001234</div>
                  </div>
                </div>
                <div class="border-t border-dashed border-gray-300 pt-6">
                  <div class="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <div class="text-gray-500">Date</div>
                      <div class="font-semibold">March 15, 2024</div>
                    </div>
                    <div>
                      <div class="text-gray-500">Time</div>
                      <div class="font-semibold">8:00 PM</div>
                    </div>
                    <div>
                      <div class="text-gray-500">Venue</div>
                      <div class="font-semibold">City Arena</div>
                    </div>
                    <div>
                      <div class="text-gray-500">Price</div>
                      <div class="font-semibold text-green-600">$89.00</div>
                    </div>
                  </div>
                </div>
                <div class="mt-6 flex justify-center">
                  <div class="w-20 h-20 bg-gray-900 rounded-lg flex items-center justify-center">
                    <i class="fa fa-qrcode text-white text-2xl"></i>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  `
  });


  addBlock('show-section', {
    label: 'Show Section',
    category: 'Ticket Landing',
    media: '<i class="fa fa-star" fill="currentColor" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
        <section class="hero-section" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 100px 0; text-align: center; color: white;">
          <div class="container">
            <h1 style="font-size: 3.5rem; font-weight: bold; margin-bottom: 1rem;">Book Your Tickets Now!</h1>
            <p style="font-size: 1.3rem; margin-bottom: 2rem; opacity: 0.9;">Don't miss out on this amazing event</p>
            <a href="#booking" class="btn btn-warning btn-lg" style="padding: 15px 40px; font-size: 1.2rem; border-radius: 50px;">Reserve Now</a>
          </div>
        </section>
      `
  });



  addBlock('hero-one', {
    label: 'Event Info',
    category: 'Ticket Landing',
    media: '<i class="fa fa-info-circle" fill="currentColor" style="font-size: 3rem; margin-right: 10px;" ></i>',
    content: `
  <section class="relative overflow-hidden">
    <div class="absolute inset-0 bg-gradient-to-br from-indigo-50 to-green-800"></div>
    <div class="relative mx-auto max-w-7xl px-6 py-20 lg:py-28 grid lg:grid-cols-2 gap-10 items-center">
      <div>
        <h1 class="text-4xl lg:text-6xl font-extrabold tracking-tight">
          Build beautiful landing pages <span class="text-indigo-600 ">fast</span>
        </h1>
        <p class="mt-5 text-lg text-gray-600">
          Drag, drop, and publish high-converting pages with our visual editor—no code required.
        </p>
        <div class="mt-8 flex items-center gap-3">
          <a href="#start" class="px-6 py-3 rounded-xl bg-indigo-600 text-white hover:bg-indigo-700">Start building</a>
          <a href="#demo" class="px-6 py-3 rounded-xl border border-gray-300 hover:bg-gray-50">View demo</a>
        </div>
        <div class="mt-6 text-sm text-gray-500">Free 14-day trial • No credit card required</div>
      </div>
      <div class="bg-white rounded-2xl shadow-xl border border-gray-100 p-4">
        <div class="aspect-video rounded-xl bg-gray-50 grid place-items-center text-gray-400">
          Editor preview / screenshot
        </div>
      </div>
    </div>
  </section>
      `
  });



  // COUNTDOWN TIMER

  addBlock('countdown-timer', {
    label: 'Countdown Timer',
    category: 'Ticket Landing',
    media: '<i class="fa fa-clock-o" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
    <section class="py-16 bg-gradient-to-r from-red-500 to-pink-500">
      <div class="mx-auto max-w-4xl px-6 text-center">
        <h2 class="text-3xl md:text-4xl font-bold text-white mb-4">
          Event Starts In
        </h2>
        <p class="text-xl text-red-100 mb-12">Don't miss out on this amazing experience!</p>
        
        <div class="grid grid-cols-4 gap-4 md:gap-8 mb-12">
          <div class="bg-white rounded-xl p-4 md:p-6 shadow-lg">
            <div class="text-3xl md:text-5xl font-bold text-gray-900" id="days">15</div>
            <div class="text-sm md:text-base text-gray-600 font-medium">DAYS</div>
          </div>
          <div class="bg-white rounded-xl p-4 md:p-6 shadow-lg">
            <div class="text-3xl md:text-5xl font-bold text-gray-900" id="hours">07</div>
            <div class="text-sm md:text-base text-gray-600 font-medium">HOURS</div>
          </div>
          <div class="bg-white rounded-xl p-4 md:p-6 shadow-lg">
            <div class="text-3xl md:text-5xl font-bold text-gray-900" id="minutes">23</div>
            <div class="text-sm md:text-base text-gray-600 font-medium">MINUTES</div>
          </div>
          <div class="bg-white rounded-xl p-4 md:p-6 shadow-lg">
            <div class="text-3xl md:text-5xl font-bold text-gray-900" id="seconds">45</div>
            <div class="text-sm md:text-base text-gray-600 font-medium">SECONDS</div>
          </div>
        </div>
        
        <a href="#tickets" class="inline-flex items-center px-8 py-4 bg-white text-red-500 font-bold text-lg rounded-xl hover:bg-gray-100 transition-colors shadow-lg">
          <i class="fa fa-ticket-alt mr-3"></i>
          Get Tickets Now
        </a>
      </div>
    </section>
  `
  });

  // SOCIAL PROOF/TESTIMONIALS

  addBlock('social-proof', {
    label: 'Social Proof',
    category: 'Ticket Landing',
    media: '<i class="fa fa-star" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
    <section class="py-20 bg-white">
      <div class="mx-auto max-w-7xl px-6">
        <div class="text-center mb-16">
          <h2 class="text-4xl font-bold text-gray-900 mb-4">What People Are Saying</h2>
          <p class="text-xl text-gray-600">Join thousands of satisfied attendees</p>
        </div>
        
        <!-- Stats -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8 mb-16">
          <div class="text-center">
            <div class="text-4xl font-bold text-blue-600 mb-2">50K+</div>
            <div class="text-gray-600">Happy Customers</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-green-600 mb-2">4.9★</div>
            <div class="text-gray-600">Average Rating</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-purple-600 mb-2">98%</div>
            <div class="text-gray-600">Would Recommend</div>
          </div>
          <div class="text-center">
            <div class="text-4xl font-bold text-yellow-600 mb-2">24/7</div>
            <div class="text-gray-600">Customer Support</div>
          </div>
        </div>
        
        <!-- Testimonials -->
        <div class="grid md:grid-cols-3 gap-8">
          <div class="bg-gray-50 rounded-2xl p-8">
            <div class="flex items-center mb-4">
              <div class="flex text-yellow-400">
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
              </div>
            </div>
            <p class="text-gray-700 mb-6 italic">"Amazing experience! The venue was perfect and the organization was flawless. Will definitely attend again next year."</p>
            <div class="flex items-center">
              <div class="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center text-white font-bold mr-4">
                S
              </div>
              <div>
                <div class="font-semibold text-gray-900">Sarah Johnson</div>
                <div class="text-gray-500 text-sm">Verified Purchase</div>
              </div>
            </div>
          </div>
          
          <div class="bg-gray-50 rounded-2xl p-8">
            <div class="flex items-center mb-4">
              <div class="flex text-yellow-400">
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
              </div>
            </div>
            <p class="text-gray-700 mb-6 italic">"Easy ticket purchase process and great customer service. The mobile tickets made entry super smooth."</p>
            <div class="flex items-center">
              <div class="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center text-white font-bold mr-4">
                M
              </div>
              <div>
                <div class="font-semibold text-gray-900">Mike Chen</div>
                <div class="text-gray-500 text-sm">Verified Purchase</div>
              </div>
            </div>
          </div>
          
          <div class="bg-gray-50 rounded-2xl p-8">
            <div class="flex items-center mb-4">
              <div class="flex text-yellow-400">
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
                <i class="fa fa-star"></i>
              </div>
            </div>
            <p class="text-gray-700 mb-6 italic">"Incredible value for money. The VIP experience exceeded all expectations. Highly recommend!"</p>
            <div class="flex items-center">
              <div class="w-12 h-12 bg-purple-500 rounded-full flex items-center justify-center text-white font-bold mr-4">
                A
              </div>
              <div>
                <div class="font-semibold text-gray-900">Amanda Rodriguez</div>
                <div class="text-gray-500 text-sm">VIP Customer</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  `
  });

  // FAQ SECTION

  addBlock('faq-section', {
    label: 'FAQ Section',
    category: 'Ticket Landing',
    media: '<i class="fa fa-question-circle" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
    <section class="py-20 bg-gray-50" x-data="{ openFaq: null }">
      <div class="mx-auto max-w-4xl px-6">
        <div class="text-center mb-16">
          <h2 class="text-4xl font-bold text-gray-900 mb-4">Frequently Asked Questions</h2>
          <p class="text-xl text-gray-600">Everything you need to know about your tickets</p>
        </div>
        
        <div class="space-y-4">
          <div class="bg-white rounded-lg border border-gray-200">
            <button class="w-full px-6 py-4 text-left flex justify-between items-center" @click="openFaq = openFaq === 1 ? null : 1">
              <span class="font-semibold text-gray-900">How do I receive my tickets?</span>
              <i class="fa fa-chevron-down transition-transform duration-200" :class="{ 'rotate-180': openFaq === 1 }"></i>
            </button>
            <div x-show="openFaq === 1" x-collapse class="px-6 pb-4">
              <p class="text-gray-600">All tickets are delivered electronically via email and SMS. You'll receive your mobile tickets within 24 hours of purchase. Simply show the QR code on your phone at the venue entrance.</p>
            </div>
          </div>
          
          <div class="bg-white rounded-lg border border-gray-200">
            <button class="w-full px-6 py-4 text-left flex justify-between items-center" @click="openFaq = openFaq === 2 ? null : 2">
              <span class="font-semibold text-gray-900">Can I refund or exchange my tickets?</span>
              <i class="fa fa-chevron-down transition-transform duration-200" :class="{ 'rotate-180': openFaq === 2 }"></i>
            </button>
            <div x-show="openFaq === 2" x-collapse class="px-6 pb-4">
              <p class="text-gray-600">Tickets can be refunded up to 7 days before the event for a full refund minus processing fees. Exchanges are available subject to availability. Contact our support team for assistance.</p>
            </div>
          </div>
          
          <div class="bg-white rounded-lg border border-gray-200">
            <button class="w-full px-6 py-4 text-left flex justify-between items-center" @click="openFaq = openFaq === 3 ? null : 3">
              <span class="font-semibold text-gray-900">What should I bring to the event?</span>
              <i class="fa fa-chevron-down transition-transform duration-200" :class="{ 'rotate-180': openFaq === 3 }"></i>
            </button>
            <div x-show="openFaq === 3" x-collapse class="px-6 pb-4">
              <p class="text-gray-600">Bring a valid photo ID and your mobile device with the ticket QR code. Small bags are allowed but subject to security checks. No outside food or drinks permitted.</p>
            </div>
          </div>
          
          <div class="bg-white rounded-lg border border-gray-200">
            <button class="w-full px-6 py-4 text-left flex justify-between items-center" @click="openFaq = openFaq === 4 ? null : 4">
              <span class="font-semibold text-gray-900">Is parking available at the venue?</span>
              <i class="fa fa-chevron-down transition-transform duration-200" :class="{ 'rotate-180': openFaq === 4 }"></i>
            </button>
            <div x-show="openFaq === 4" x-collapse class="px-6 pb-4">
              <p class="text-gray-600">Yes, venue parking is available for $15. We also recommend using public transportation or rideshare services. Street parking may be limited.</p>
            </div>
          </div>
          
          <div class="bg-white rounded-lg border border-gray-200">
            <button class="w-full px-6 py-4 text-left flex justify-between items-center" @click="openFaq = openFaq === 5 ? null : 5">
              <span class="font-semibold text-gray-900">Are there age restrictions?</span>
              <i class="fa fa-chevron-down transition-transform duration-200" :class="{ 'rotate-180': openFaq === 5 }"></i>
            </button>
            <div x-show="openFaq === 5" x-collapse class="px-6 pb-4">
              <p class="text-gray-600">This is an all-ages event. Children under 12 must be accompanied by an adult. Children under 2 do not require a ticket if sitting on a parent's lap.</p>
            </div>
          </div>
        </div>
        
        <div class="text-center mt-12">
          <p class="text-gray-600 mb-4">Still have questions?</p>
          <a href="#contact" class="inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
            <i class="fa fa-envelope mr-2"></i>
            Contact Support
          </a>
        </div>
      </div>
    </section>
  `
  });

  // FOOTER

  addBlock('ticket-footer', {
    label: 'Footer',
    category: 'Ticket Landing',
    media: '<i class="fa fa-link" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
    <footer class="bg-gray-900 text-white py-16">
      <div class="mx-auto max-w-7xl px-6">
        <div class="grid md:grid-cols-4 gap-8">
          <div>
            <h3 class="text-xl font-bold mb-4">EventTickets</h3>
            <p class="text-gray-400 mb-4">Your trusted source for event tickets and unforgettable experiences.</p>
            <div class="flex space-x-4">
              <a href="#" class="text-gray-400 hover:text-white">
                <i class="fab fa-facebook text-xl"></i>
              </a>
              <a href="#" class="text-gray-400 hover:text-white">
                <i class="fab fa-twitter text-xl"></i>
              </a>
              <a href="#" class="text-gray-400 hover:text-white">
                <i class="fab fa-instagram text-xl"></i>
              </a>
            </div>
          </div>
          
          <div>
            <h4 class="font-semibold mb-4">Quick Links</h4>
            <ul class="space-y-2 text-gray-400">
              <li><a href="#" class="hover:text-white">Buy Tickets</a></li>
              <li><a href="#" class="hover:text-white">Event Info</a></li>
              <li><a href="#" class="hover:text-white">Venue Details</a></li>
              <li><a href="#" class="hover:text-white">Parking Info</a></li>
            </ul>
          </div>
          
          <div>
            <h4 class="font-semibold mb-4">Support</h4>
            <ul class="space-y-2 text-gray-400">
              <li><a href="#" class="hover:text-white">Contact Us</a></li>
              <li><a href="#" class="hover:text-white">FAQ</a></li>
              <li><a href="#" class="hover:text-white">Refund Policy</a></li>
              <li><a href="#" class="hover:text-white">Terms of Service</a></li>
            </ul>
          </div>
          
          <div>
            <h4 class="font-semibold mb-4">Contact Info</h4>
            <div class="space-y-2 text-gray-400">
              <div class="flex items-center">
                <i class="fa fa-phone mr-2"></i>
                <span>(555) 123-4567</span>
              </div>
              <div class="flex items-center">
                <i class="fa fa-envelope mr-2"></i>
                <span>support@eventtickets.com</span>
              </div>
              <div class="flex items-center">
                <i class="fa fa-clock mr-2"></i>
                <span>24/7 Support Available</span>
              </div>
            </div>
          </div>
        </div>
        
        <div class="border-t border-gray-800 mt-12 pt-8 text-center text-gray-400">
          <p>&copy; 2024 EventTickets. All rights reserved. | Privacy Policy | Terms of Service</p>
        </div>
      </div>
    </footer>
  `
  });

  // NEWSLETTER SIGNUP

  addBlock('newsletter-signup', {
    label: 'Newsletter Signup',
    category: 'Ticket Landing',
    media: '<i class="fa fa-envelope" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
    <section class="py-20 bg-gradient-to-r from-blue-600 to-purple-600">
      <div class="mx-auto max-w-4xl px-6 text-center">
        <h2 class="text-4xl font-bold text-white mb-4">Stay Updated</h2>
        <p class="text-xl text-blue-100 mb-8">Get notified about upcoming events and exclusive offers</p>
        
        <div class="max-w-lg mx-auto">
          <div class="flex flex-col sm:flex-row gap-4">
            <input type="email" placeholder="Enter your email address" 
                   class="flex-1 bg-white px-6 py-4 rounded-lg border-0 focus:ring-2 focus:ring-white focus:outline-none text-gray-900">
            <button class="px-8 py-4 bg-yellow-400 text-black font-semibold rounded-lg hover:bg-yellow-300 transition-colors">
              Subscribe
            </button>
          </div>
          <p class="text-sm text-blue-200 mt-4">
            We respect your privacy. Unsubscribe at any time.
          </p>
        </div>
        
        <div class="mt-12 grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
          <div class="text-blue-100">
            <i class="fa fa-bell text-3xl mb-3"></i>
            <h3 class="font-semibold mb-2">Early Access</h3>
            <p class="text-sm">Get tickets before general sale</p>
          </div>
          <div class="text-blue-100">
            <i class="fa fa-percent text-3xl mb-3"></i>
            <h3 class="font-semibold mb-2">Exclusive Discounts</h3>
            <p class="text-sm">Subscriber-only special offers</p>
          </div>
          <div class="text-blue-100">
            <i class="fa fa-calendar text-3xl mb-3"></i>
            <h3 class="font-semibold mb-2">Event Updates</h3>
            <p class="text-sm">Never miss an upcoming event</p>
          </div>
        </div>
      </div>
    </section>
  `
  });



  // TICKET RESERVATION PAGE BUILDER COMPONENTS

  // HERO SECTIONS
  addBlock('ticket-hero-modern', {
    label: 'Modern Ticket Hero',
    category: 'Ticket Landing',
    media: '<i class="fa fa-ticket" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
  <section class="relative overflow-hidden bg-gradient-to-br from-purple-900 via-blue-900 to-indigo-900 min-h-screen flex items-center" x-data="{ scrollY: 0 }" x-init="window.addEventListener('scroll', () => scrollY = window.scrollY)">
    <div class="absolute inset-0 opacity-10" :style="'transform: translateY(' + scrollY * 0.5 + 'px)'">
      <div class="absolute inset-0" style="background-image: url('data:image/svg+xml,<svg width="60" height="60" viewBox="0 0 60 60" xmlns="http://www.w3.org/2000/svg"><g fill="none" fill-rule="evenodd"><g fill="%23ffffff" fill-opacity="0.1"><circle cx="30" cy="30" r="2"/></g></g></svg>'); background-size: 60px 60px;"></div>
    </div>
    
    <!-- Floating shapes -->
    <div class="absolute top-20 left-10 w-20 h-20 bg-yellow-400/20 rounded-full blur-xl animate-pulse"></div>
    <div class="absolute bottom-20 right-10 w-32 h-32 bg-purple-400/20 rounded-full blur-xl animate-pulse delay-1000"></div>
    
    <div class="relative mx-auto max-w-7xl px-6 py-16">
      <div class="grid lg:grid-cols-2 gap-12 items-center">
        <div class="text-center lg:text-left" x-data="{ visible: false }" x-init="setTimeout(() => visible = true, 300)">
          <div class="inline-flex items-center px-4 py-2 rounded-full bg-gradient-to-r from-yellow-400 to-orange-500 text-black text-sm font-semibold mb-6 transform transition-all duration-1000" :class="visible ? 'translate-y-0 opacity-100' : 'translate-y-10 opacity-0'">
            <i class="fa fa-fire mr-2"></i>
            🔥 Selling Fast
          </div>
          <h1 class="text-4xl lg:text-7xl font-black tracking-tight text-white leading-tight transform transition-all duration-1000 delay-200" :class="visible ? 'translate-y-0 opacity-100' : 'translate-y-10 opacity-0'">
            Experience
            <span class="text-transparent bg-clip-text bg-gradient-to-r from-yellow-400 via-pink-500 to-purple-500 animate-pulse">Magic</span>
            <br>Live in Concert
          </h1>
          <p class="mt-6 text-xl text-gray-300 leading-relaxed transform transition-all duration-1000 delay-400" :class="visible ? 'translate-y-0 opacity-100' : 'translate-y-10 opacity-0'">
            Join thousands of fans for an unforgettable night of music, lights, and pure energy. This is more than a concert—it's a movement.
          </p>
          
          <!-- Event Info Cards -->
          <div class="mt-8 grid grid-cols-1 sm:grid-cols-3 gap-4" x-data="{ showCards: false }" x-init="setTimeout(() => showCards = true, 800)">
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-4 transform transition-all duration-500" :class="showCards ? 'translate-y-0 opacity-100' : 'translate-y-10 opacity-0'">
              <i class="fa fa-calendar-alt text-yellow-400 text-2xl mb-2"></i>
              <div class="font-semibold text-white">March 15</div>
              <div class="text-sm text-gray-300">2024</div>
            </div>
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-4 transform transition-all duration-500 delay-100" :class="showCards ? 'translate-y-0 opacity-100' : 'translate-y-10 opacity-0'">
              <i class="fa fa-clock text-yellow-400 text-2xl mb-2"></i>
              <div class="font-semibold text-white">8:00 PM</div>
              <div class="text-sm text-gray-300">Doors at 7:30</div>
            </div>
            <div class="bg-white/10 backdrop-blur-lg rounded-xl p-4 transform transition-all duration-500 delay-200" :class="showCards ? 'translate-y-0 opacity-100' : 'translate-y-10 opacity-0'">
              <i class="fa fa-map-marker-alt text-yellow-400 text-2xl mb-2"></i>
              <div class="font-semibold text-white">Madison Square</div>
              <div class="text-sm text-gray-300">New York</div>
            </div>
          </div>
          
          <div class="mt-10 flex flex-col sm:flex-row items-center justify-center lg:justify-start gap-4">
            <button x-data="{ hover: false }" @mouseenter="hover = true" @mouseleave="hover = false" class="group w-full sm:w-auto px-8 py-4 rounded-xl bg-gradient-to-r from-yellow-400 to-orange-500 text-black font-bold text-lg transition-all duration-300 shadow-2xl relative overflow-hidden" :class="hover ? 'scale-105' : ''">
              <div class="absolute inset-0 bg-white/20 transform scale-x-0 group-hover:scale-x-100 transition-transform duration-300 origin-left"></div>
              <span class="relative flex items-center justify-center">
                <i class="fa fa-ticket-alt mr-2"></i>Get Tickets Now
              </span>
            </button>
            <button class="w-full sm:w-auto px-8 py-4 rounded-xl border-2 border-white/30 text-white font-semibold hover:bg-white/10 transition-all duration-300 backdrop-blur-sm">
              <i class="fa fa-play mr-2"></i>Watch Trailer
            </button>
          </div>
        </div>
        
        <!-- Ticket Preview -->
        <div class="flex justify-center lg:justify-end" x-data="{ ticketVisible: false }" x-init="setTimeout(() => ticketVisible = true, 1000)">
          <div class="relative transform transition-all duration-1000" :class="ticketVisible ? 'translate-y-0 opacity-100 rotate-3' : 'translate-y-20 opacity-0 rotate-12'">
            <div class="absolute -inset-4 bg-gradient-to-r from-yellow-400 to-pink-500 rounded-2xl blur-xl opacity-50"></div>
            <div class="relative bg-white rounded-2xl shadow-2xl border border-gray-200 p-8 max-w-sm">
              <div class="flex justify-between items-start mb-6">
                <div>
                  <div class="text-2xl font-bold text-gray-900">ADMIT ONE</div>
                  <div class="text-purple-600 font-semibold">VIP Access</div>
                </div>
                <div class="text-right">
                  <div class="text-sm text-gray-500">Ticket #</div>
                  <div class="font-mono font-bold text-gray-900">VIP001</div>
                </div>
              </div>
              <div class="border-t-2 border-dashed border-gray-300 pt-6">
                <div class="space-y-3 text-sm">
                  <div class="flex justify-between">
                    <span class="text-gray-500">Artist</span>
                    <span class="font-semibold">The Midnight Stars</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-gray-500">Venue</span>
                    <span class="font-semibold">Madison Square</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-gray-500">Date & Time</span>
                    <span class="font-semibold">Mar 15, 8:00 PM</span>
                  </div>
                  <div class="flex justify-between">
                    <span class="text-gray-500">Price</span>
                    <span class="font-semibold text-green-600 text-lg">$149.00</span>
                  </div>
                </div>
              </div>
              <div class="mt-6 flex justify-center">
                <div class="w-16 h-16 bg-gradient-to-br from-purple-600 to-blue-600 rounded-lg flex items-center justify-center">
                  <i class="fa fa-qrcode text-white text-xl"></i>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
  `
  });

  // TICKET SELECTION SECTION
  addBlock('ticket-selector', {
    label: 'Ticket Selector',
    category: 'Ticket Landing',
    media: '<i class="fa fa-list" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
  <section class="py-16 bg-gray-50" x-data="ticketSelector()">
    <div class="max-w-6xl mx-auto px-6">
      <div class="text-center mb-12">
        <h2 class="text-4xl font-bold text-gray-900 mb-4">Choose Your Experience</h2>
        <p class="text-xl text-gray-600">Select the perfect tickets for an unforgettable night</p>
      </div>
      
      <div class="grid md:grid-cols-3 gap-8">
        <!-- General Admission -->
        <div class="bg-white rounded-2xl shadow-lg border border-gray-200 overflow-hidden transform transition-all duration-300 hover:scale-105 hover:shadow-2xl" :class="selectedTier === 'general' ? 'ring-4 ring-blue-500' : ''" @click="selectTier('general')">
          <div class="p-8">
            <div class="text-center mb-6">
              <div class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <i class="fa fa-users text-blue-600 text-2xl"></i>
              </div>
              <h3 class="text-2xl font-bold text-gray-900">General Admission</h3>
              <p class="text-gray-600 mt-2">Standing room, full access</p>
            </div>
            <div class="text-center mb-6">
              <div class="text-4xl font-bold text-gray-900">$89</div>
              <div class="text-gray-500">per ticket</div>
            </div>
            <ul class="space-y-3 mb-8">
              <li class="flex items-center text-gray-700">
                <i class="fa fa-check text-green-500 mr-3"></i>
                Standing room access
              </li>
              <li class="flex items-center text-gray-700">
                <i class="fa fa-check text-green-500 mr-3"></i>
                Full concert experience
              </li>
              <li class="flex items-center text-gray-700">
                <i class="fa fa-check text-green-500 mr-3"></i>
                Mobile ticket delivery
              </li>
            </ul>
            <div class="flex items-center justify-between mb-4">
              <label class="text-gray-700 font-medium">Quantity:</label>
              <div class="flex items-center space-x-3">
                <button @click="decreaseQuantity('general')" class="w-8 h-8 rounded-full bg-gray-200 hover:bg-gray-300 flex items-center justify-center">
                  <i class="fa fa-minus text-sm"></i>
                </button>
                <span class="w-8 text-center font-semibold" x-text="quantities.general"></span>
                <button @click="increaseQuantity('general')" class="w-8 h-8 rounded-full bg-gray-200 hover:bg-gray-300 flex items-center justify-center">
                  <i class="fa fa-plus text-sm"></i>
                </button>
              </div>
            </div>
            <button class="w-full py-3 rounded-lg font-semibold transition-all duration-300" :class="selectedTier === 'general' ? 'bg-blue-600 text-white hover:bg-blue-700' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'">
              <span x-show="selectedTier !== 'general'">Select Tickets</span>
              <span x-show="selectedTier === 'general'">✓ Selected</span>
            </button>
          </div>
        </div>

        <!-- VIP -->
        <div class="bg-white rounded-2xl shadow-lg border-2 border-gradient-to-r from-yellow-400 to-orange-500 overflow-hidden transform transition-all duration-300 hover:scale-105 hover:shadow-2xl relative" :class="selectedTier === 'vip' ? 'ring-4 ring-yellow-500' : ''" @click="selectTier('vip')">
          <div class="absolute top-4 right-4 bg-gradient-to-r from-yellow-400 to-orange-500 text-black px-3 py-1 rounded-full text-sm font-bold">
            POPULAR
          </div>
          <div class="p-8">
            <div class="text-center mb-6">
              <div class="w-16 h-16 bg-gradient-to-r from-yellow-100 to-orange-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <i class="fa fa-crown text-yellow-600 text-2xl"></i>
              </div>
              <h3 class="text-2xl font-bold text-gray-900">VIP Experience</h3>
              <p class="text-gray-600 mt-2">Premium seating & perks</p>
            </div>
            <div class="text-center mb-6">
              <div class="text-4xl font-bold text-gray-900">$199</div>
              <div class="text-gray-500">per ticket</div>
            </div>
            <ul class="space-y-3 mb-8">
              <li class="flex items-center text-gray-700">
                <i class="fa fa-check text-green-500 mr-3"></i>
                Reserved seating
              </li>
              <li class="flex items-center text-gray-700">
                <i class="fa fa-check text-green-500 mr-3"></i>
                VIP lounge access
              </li>
              <li class="flex items-center text-gray-700">
                <i class="fa fa-check text-green-500 mr-3"></i>
                Complimentary drinks
              </li>
              <li class="flex items-center text-gray-700">
                <i class="fa fa-check text-green-500 mr-3"></i>
                Meet & greet opportunity
              </li>
            </ul>
            <div class="flex items-center justify-between mb-4">
              <label class="text-gray-700 font-medium">Quantity:</label>
              <div class="flex items-center space-x-3">
                <button @click="decreaseQuantity('vip')" class="w-8 h-8 rounded-full bg-gray-200 hover:bg-gray-300 flex items-center justify-center">
                  <i class="fa fa-minus text-sm"></i>
                </button>
                <span class="w-8 text-center font-semibold" x-text="quantities.vip"></span>
                <button @click="increaseQuantity('vip')" class="w-8 h-8 rounded-full bg-gray-200 hover:bg-gray-300 flex items-center justify-center">
                  <i class="fa fa-plus text-sm"></i>
                </button>
              </div>
            </div>
            <button class="w-full py-3 rounded-lg font-semibold transition-all duration-300" :class="selectedTier === 'vip' ? 'bg-gradient-to-r from-yellow-400 to-orange-500 text-black hover:from-yellow-300 hover:to-orange-400' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'">
              <span x-show="selectedTier !== 'vip'">Select Tickets</span>
              <span x-show="selectedTier === 'vip'">✓ Selected</span>
            </button>
          </div>
        </div>

        <!-- Premium -->
        <div class="bg-gradient-to-br from-purple-900 to-indigo-900 rounded-2xl shadow-lg overflow-hidden transform transition-all duration-300 hover:scale-105 hover:shadow-2xl text-white relative" :class="selectedTier === 'premium' ? 'ring-4 ring-purple-400' : ''" @click="selectTier('premium')">
          <div class="absolute top-4 right-4 bg-purple-400 text-black px-3 py-1 rounded-full text-sm font-bold">
            EXCLUSIVE
          </div>
          <div class="p-8">
            <div class="text-center mb-6">
              <div class="w-16 h-16 bg-purple-400/20 rounded-full flex items-center justify-center mx-auto mb-4">
                <i class="fa fa-gem text-purple-300 text-2xl"></i>
              </div>
              <h3 class="text-2xl font-bold">Premium Package</h3>
              <p class="text-purple-200 mt-2">Ultimate luxury experience</p>
            </div>
            <div class="text-center mb-6">
              <div class="text-4xl font-bold">$399</div>
              <div class="text-purple-200">per ticket</div>
            </div>
            <ul class="space-y-3 mb-8">
              <li class="flex items-center">
                <i class="fa fa-check text-green-400 mr-3"></i>
                Front row seating
              </li>
              <li class="flex items-center">
                <i class="fa fa-check text-green-400 mr-3"></i>
                Private hospitality suite
              </li>
              <li class="flex items-center">
                <i class="fa fa-check text-green-400 mr-3"></i>
                Backstage tour
              </li>
              <li class="flex items-center">
                <i class="fa fa-check text-green-400 mr-3"></i>
                Exclusive merchandise
              </li>
              <li class="flex items-center">
                <i class="fa fa-check text-green-400 mr-3"></i>
                Professional photo
              </li>
            </ul>
            <div class="flex items-center justify-between mb-4">
              <label class="font-medium">Quantity:</label>
              <div class="flex items-center space-x-3">
                <button @click="decreaseQuantity('premium')" class="w-8 h-8 rounded-full bg-white/20 hover:bg-white/30 flex items-center justify-center">
                  <i class="fa fa-minus text-sm"></i>
                </button>
                <span class="w-8 text-center font-semibold" x-text="quantities.premium"></span>
                <button @click="increaseQuantity('premium')" class="w-8 h-8 rounded-full bg-white/20 hover:bg-white/30 flex items-center justify-center">
                  <i class="fa fa-plus text-sm"></i>
                </button>
              </div>
            </div>
            <button class="w-full py-3 rounded-lg font-semibold transition-all duration-300" :class="selectedTier === 'premium' ? 'bg-purple-400 text-black hover:bg-purple-300' : 'bg-white/20 hover:bg-white/30'">
              <span x-show="selectedTier !== 'premium'">Select Tickets</span>
              <span x-show="selectedTier === 'premium'">✓ Selected</span>
            </button>
          </div>
        </div>
      </div>
      
      <!-- Total and Checkout -->
      <div class="mt-12 max-w-md mx-auto" x-show="totalTickets > 0" x-transition>
        <div class="bg-white rounded-xl shadow-lg p-6 border border-gray-200">
          <div class="flex justify-between items-center mb-4">
            <span class="text-lg font-semibold text-gray-900">Total Tickets:</span>
            <span class="text-lg font-bold text-blue-600" x-text="totalTickets"></span>
          </div>
          <div class="flex justify-between items-center mb-6">
            <span class="text-xl font-bold text-gray-900">Total Price:</span>
            <span class="text-2xl font-bold text-green-600" x-text="'$' + totalPrice.toFixed(2)"></span>
          </div>
          <button class="w-full py-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-bold text-lg rounded-xl hover:from-blue-700 hover:to-purple-700 transition-all duration-300 transform hover:scale-105">
            <i class="fa fa-credit-card mr-2"></i>
            Proceed to Checkout
          </button>
        </div>
      </div>
    </div>
    
    <script>
      function ticketSelector() {
        return {
          selectedTier: null,
          quantities: {
            general: 0,
            vip: 0,
            premium: 0
          },
          prices: {
            general: 89,
            vip: 199,
            premium: 399
          },
          
          get totalTickets() {
            return this.quantities.general + this.quantities.vip + this.quantities.premium;
          },
          
          get totalPrice() {
            return (this.quantities.general * this.prices.general) + 
                   (this.quantities.vip * this.prices.vip) + 
                   (this.quantities.premium * this.prices.premium);
          },
          
          selectTier(tier) {
            this.selectedTier = tier;
            if (this.quantities[tier] === 0) {
              this.quantities[tier] = 1;
            }
          },
          
          increaseQuantity(tier) {
            if (this.quantities[tier] < 8) {
              this.quantities[tier]++;
              this.selectedTier = tier;
            }
          },
          
          decreaseQuantity(tier) {
            if (this.quantities[tier] > 0) {
              this.quantities[tier]--;
            }
          }
        }
      }
    </script>
  </section>
  `
  });

  // COUNTDOWN TIMER SECTION
  addBlock('ticket-countdown', {
    label: 'Countdown Timer',
    category: 'Ticket Landing',
    media: '<i class="fa fa-clock" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
  <section class="py-16 bg-gradient-to-r from-red-600 to-pink-600 text-white" x-data="countdownTimer()">
    <div class="max-w-4xl mx-auto px-6 text-center">
      <div class="mb-8">
        <h2 class="text-3xl lg:text-5xl font-bold mb-4">⚡ Flash Sale Ending Soon!</h2>
        <p class="text-xl text-red-100">Get 25% off all tickets - Limited time offer</p>
      </div>
      
      <div class="grid grid-cols-4 gap-4 max-w-2xl mx-auto mb-8">
        <div class="bg-white/10 backdrop-blur-lg rounded-xl p-4">
          <div class="text-3xl lg:text-4xl font-bold" x-text="timeLeft.days"></div>
          <div class="text-sm text-red-100">Days</div>
        </div>
        <div class="bg-white/10 backdrop-blur-lg rounded-xl p-4">
          <div class="text-3xl lg:text-4xl font-bold" x-text="timeLeft.hours"></div>
          <div class="text-sm text-red-100">Hours</div>
        </div>
        <div class="bg-white/10 backdrop-blur-lg rounded-xl p-4">
          <div class="text-3xl lg:text-4xl font-bold" x-text="timeLeft.minutes"></div>
          <div class="text-sm text-red-100">Minutes</div>
        </div>
        <div class="bg-white/10 backdrop-blur-lg rounded-xl p-4">
          <div class="text-3xl lg:text-4xl font-bold" x-text="timeLeft.seconds"></div>
          <div class="text-sm text-red-100">Seconds</div>
        </div>
      </div>
      
      <div class="flex justify-center">
        <button class="px-10 py-4 bg-white text-red-600 font-bold text-lg rounded-xl hover:bg-gray-100 transition-all duration-300 transform hover:scale-105 shadow-xl">
          <i class="fa fa-bolt mr-2"></i>
          Claim Your Discount Now
        </button>
      </div>
    </div>
    
    <script>
      function countdownTimer() {
        return {
          timeLeft: {
            days: 0,
            hours: 0,
            minutes: 0,
            seconds: 0
          },
          
          init() {
            this.updateCountdown();
            setInterval(() => {
              this.updateCountdown();
            }, 1000);
          },
          
          updateCountdown() {
            const now = new Date().getTime();
            const eventDate = new Date('2024-03-10 23:59:59').getTime();
            const distance = eventDate - now;
            
            if (distance > 0) {
              this.timeLeft.days = Math.floor(distance / (1000 * 60 * 60 * 24));
              this.timeLeft.hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
              this.timeLeft.minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
              this.timeLeft.seconds = Math.floor((distance % (1000 * 60)) / 1000);
            } else {
              this.timeLeft = { days: 0, hours: 0, minutes: 0, seconds: 0 };
            }
          }
        }
      }
    </script>
  </section>
  `
  });


  // VENUE SEATING MAP
  addBlock('venue-seating', {
    label: 'Venue Seating Map',
    category: 'Ticket Landing',
    media: '<i class="fa fa-map" style="font-size: 3rem; margin-right: 10px;"></i>',
    content: `
  <section class="py-16 bg-white" x-data="seatingMap()">
    <div class="max-w-6xl mx-auto px-6">
      <div class="text-center mb-12">
        <h2 class="text-4xl font-bold text-gray-900 mb-4">Venue Layout</h2>
        <p class="text-xl text-gray-600">Explore the venue and choose your preferred section</p>
      </div>
      
      <div class="grid lg:grid-cols-3 gap-8">
        <!-- Seating Map -->
        <div class="lg:col-span-2">
          <div class="bg-gray-50 rounded-2xl p-8 border border-gray-200">
            <div class="text-center mb-8">
              <div class="inline-block bg-gradient-to-r from-purple-600 to-pink-600 text-white px-6 py-3 rounded-lg font-bold text-lg">
                🎤 STAGE
              </div>
            </div>
            
            <!-- Seating Areas -->
            <div class="space-y-6">
              <!-- Premium Front Row -->
              <div class="text-center">
                <div class="grid grid-cols-10 gap-1 max-w-md mx-auto">
                  <template x-for="seat in 10">
                    <button @click="selectSeat('premium', seat)" class="w-8 h-8 rounded-md border-2 transition-all duration-200" :class="selectedSeats.premium.includes(seat) ? 'bg-purple-600 border-purple-600 text-white' : 'bg-purple-100 border-purple-300 hover:bg-purple-200'">
                      <i class="fa fa-chair text-xs"></i>
                    </button>
                  </template>
                </div>
                <div class="text-sm text-gray-600 mt-2">Premium - $399</div>
              </div>
              
              <!-- VIP Section -->
              <div class="text-center">
                <div class="grid grid-cols-12 gap-1 max-w-lg mx-auto">
                  <template x-for="seat in 24">
                    <button @click="selectSeat('vip', seat)" class="w-7 h-7 rounded-md border-2 transition-all duration-200" :class="selectedSeats.vip.includes(seat) ? 'bg-yellow-500 border-yellow-500 text-black' : 'bg-yellow-100 border-yellow-300 hover:bg-yellow-200'">
                      <i class="fa fa-chair text-xs"></i>
                    </button>
                  </template>
                </div>
                <div class="text-sm text-gray-600 mt-2">VIP Seating - $199</div>
              </div>
              
              <!-- General Standing -->
              <div class="text-center">
                <div class="bg-blue-100 border-2 border-blue-300 rounded-lg p-8 max-w-lg mx-auto">
                  <div class="grid grid-cols-8 gap-2">
                    <template x-for="area in 32">
                      <div class="w-4 h-4 bg-blue-300 rounded-full" :class="generalSelected ? 'bg-blue-600' : ''"></div>
                    </template>
                  </div>
                  <div class="mt-4">
                    <button @click="toggleGeneral()" class="px-6 py-2 rounded-lg font-semibold transition-all duration-300" :class="generalSelected ? 'bg-blue-600 text-white' : 'bg-blue-200 text-blue-800 hover:bg-blue-300'">
                      <span x-show="!generalSelected">Select General Admission</span>
                      <span x-show="generalSelected">✓ General Selected</span>
                    </button>
                  </div>
                </div>
                <div class="text-sm text-gray-600 mt-2">General Admission - $89</div>
              </div>
            </div>
            
            <!-- Legend -->
            <div class="flex justify-center mt-8 space-x-6 text-sm">
              <div class="flex items-center">
                <div class="w-4 h-4 bg-gray-300 rounded-md mr-2"></div>
                <span>Available</span>
              </div>
              <div class="flex items-center">
                <div class="w-4 h-4 bg-blue-600 rounded-md mr-2"></div>
                <span>Selected</span>
              </div>
              <div class="flex items-center">
                <div class="w-4 h-4 bg-red-400 rounded-md mr-2"></div>
                <span>Sold Out</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Selection Summary -->
        <div class="lg:col-span-1">
          <div class="bg-gray-900 text-white rounded-2xl p-6 sticky top-6">
            <h3 class="text-xl font-bold mb-6">Your Selection</h3>
            
            <div class="space-y-4" x-show="hasSelection">
              <div x-show="selectedSeats.premium.length > 0">
                <div class="flex justify-between items-center py-2 border-b border-gray-700">
                  <span>Premium Seats</span>
                  <span x-text="selectedSeats.premium.length + ' × $399'"></span>
                </div>
              </div>
              <div x-show="selectedSeats.vip.length > 0">
                <div class="flex justify-between items-center py-2 border-b border-gray-700">
                  <span>VIP Seats</span>
                  <span x-text="selectedSeats.vip.length + ' × $199'"></span>
                </div>
              </div>
              <div x-show="generalSelected">
                <div class="flex justify-between items-center py-2 border-b border-gray-700">
                  <span>General Admission</span>
                  <span>1 × $89</span>
                </div>
              </div>
              
              <div class="pt-4">
                <div class="flex justify-between items-center text-lg font-bold">
                  <span>Total:</span>
                  <span x-text="' + totalAmount.toFixed(2)"></span>
                </div>
              </div>
              
              <button class="w-full mt-6 py-3 bg-gradient-to-r from-green-500 to-emerald-600 text-white font-bold rounded-lg hover:from-green-600 hover:to-emerald-700 transition-all duration-300">
                <i class="fa fa-lock mr-2"></i>
                Secure Checkout
              </button>
            </div>
            
            <div x-show="!hasSelection" class="text-center text-gray-400">
              <i class="fa fa-mouse-pointer text-3xl mb-4 opacity-50"></i>
              <p>Select seats from the map to see your total</p>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <script>
      function seatingMap() {
        return {
          selectedSeats: {
            premium: [],
            vip: [],
            general: []
          },
          generalSelected: false,
          
          get hasSelection() {
            return this.selectedSeats.premium.length > 0 || 
                   this.selectedSeats.vip.length > 0 || 
                   this.generalSelected;
          },
          
          get totalAmount() {
            return (this.selectedSeats.premium.length * 399) + 
                   (this.selectedSeats.vip.length * 199) + 
                   (this.generalSelected ? 89 : 0);
          },
          
          selectSeat(tier, seatNumber) {
            const index = this.selectedSeats[tier].indexOf(seatNumber);
            if (index > -1) {
              this.selectedSeats[tier].splice(index, 1);
            } else {
              if (this.selectedSeats[tier].length < 6) {
                this.selectedSeats[tier].push(seatNumber);
              }
            }
          },
          
          toggleGeneral() {
            this.generalSelected = !this.generalSelected;
          }
        }
      }
    </script>
  </section>
  `});
};


